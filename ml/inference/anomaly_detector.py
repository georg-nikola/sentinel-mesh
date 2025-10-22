"""
Anomaly Detection Module for Sentinel Mesh
Uses multiple ML algorithms to detect anomalies in Kubernetes metrics
"""

import asyncio
import logging
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional, Tuple
import json

# ML imports
from sklearn.ensemble import IsolationForest
from sklearn.cluster import DBSCAN
from sklearn.preprocessing import StandardScaler
from sklearn.decomposition import PCA
import tensorflow as tf
from pyod.models.auto_encoder import AutoEncoder
from pyod.models.lof import LOF

# Internal imports
from core.data_processor import DataProcessor
from core.model_manager import ModelManager
from models.anomaly_models import AnomalyModel, AnomalyResult


class AnomalyDetector:
    """Main anomaly detection class that orchestrates multiple detection algorithms."""

    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.logger = logging.getLogger(__name__)

        # Initialize components
        self.data_processor = DataProcessor(config)
        self.model_manager = ModelManager(config, "anomaly")

        # Detection algorithms
        self.algorithms = {
            "isolation_forest": None,
            "autoencoder": None,
            "dbscan": None,
            "lof": None,
            "statistical": None,
        }

        # Model parameters
        self.contamination = config.get("anomaly_detection", {}).get(
            "contamination", 0.1
        )
        self.threshold = config.get("anomaly_detection", {}).get("threshold", 0.8)
        self.window_size = config.get("anomaly_detection", {}).get("window_size", 100)

        # State tracking
        self.is_trained = False
        self.last_training = None
        self.feature_scaler = StandardScaler()
        self.pca_transformer = PCA(n_components=0.95)  # Keep 95% of variance

        # Initialize models
        self._initialize_models()

    def _initialize_models(self):
        """Initialize all anomaly detection models."""
        try:
            # Isolation Forest
            self.algorithms["isolation_forest"] = IsolationForest(
                contamination=self.contamination, random_state=42, n_jobs=-1
            )

            # DBSCAN for clustering-based anomaly detection
            self.algorithms["dbscan"] = DBSCAN(eps=0.5, min_samples=5)

            # Local Outlier Factor
            self.algorithms["lof"] = LOF(contamination=self.contamination)

            # AutoEncoder (will be initialized after seeing data shape)
            self.algorithms["autoencoder"] = None

            self.logger.info("Anomaly detection models initialized")

        except Exception as e:
            self.logger.error(f"Error initializing models: {e}")
            raise

    async def detect_anomalies(self, metrics_batch: List[Dict]) -> List[AnomalyResult]:
        """Detect anomalies in a batch of metrics."""
        if not metrics_batch:
            return []

        try:
            # Process metrics into feature matrix
            features_df = await self.data_processor.process_metrics_batch(metrics_batch)

            if features_df.empty:
                return []

            # Ensure models are trained
            if not self.is_trained:
                await self._train_models(features_df)

            # Detect anomalies using all algorithms
            anomalies = await self._detect_with_ensemble(features_df, metrics_batch)

            self.logger.debug(
                f"Detected {len(anomalies)} anomalies from {len(metrics_batch)} metrics"
            )
            return anomalies

        except Exception as e:
            self.logger.error(f"Error detecting anomalies: {e}")
            return []

    async def _detect_with_ensemble(
        self, features_df: pd.DataFrame, original_metrics: List[Dict]
    ) -> List[AnomalyResult]:
        """Use ensemble of algorithms to detect anomalies."""
        anomalies = []

        try:
            # Prepare features
            features = self._prepare_features(features_df)

            if features is None or len(features) == 0:
                return []

            # Get predictions from each algorithm
            predictions = {}
            scores = {}

            # Isolation Forest
            if self.algorithms["isolation_forest"] is not None:
                pred = self.algorithms["isolation_forest"].predict(features)
                score = self.algorithms["isolation_forest"].decision_function(features)
                predictions["isolation_forest"] = pred
                scores["isolation_forest"] = score

            # LOF
            if self.algorithms["lof"] is not None:
                pred = self.algorithms["lof"].predict(features)
                score = self.algorithms["lof"].decision_function(features)
                predictions["lof"] = pred
                scores["lof"] = score

            # AutoEncoder
            if self.algorithms["autoencoder"] is not None:
                pred = self.algorithms["autoencoder"].predict(features)
                score = self.algorithms["autoencoder"].decision_function(features)
                predictions["autoencoder"] = pred
                scores["autoencoder"] = score

            # Statistical anomaly detection
            statistical_anomalies = self._detect_statistical_anomalies(features_df)

            # Combine results using voting
            ensemble_anomalies = self._combine_predictions(
                predictions,
                scores,
                features_df,
                original_metrics,
                statistical_anomalies,
            )

            return ensemble_anomalies

        except Exception as e:
            self.logger.error(f"Error in ensemble detection: {e}")
            return []

    def _prepare_features(self, features_df: pd.DataFrame) -> Optional[np.ndarray]:
        """Prepare features for anomaly detection."""
        try:
            if features_df.empty:
                return None

            # Select numeric columns
            numeric_columns = features_df.select_dtypes(include=[np.number]).columns
            if len(numeric_columns) == 0:
                return None

            features = features_df[numeric_columns].fillna(0)

            # Scale features
            if hasattr(self.feature_scaler, "scale_"):
                features_scaled = self.feature_scaler.transform(features)
            else:
                features_scaled = self.feature_scaler.fit_transform(features)

            # Apply PCA if needed
            if features_scaled.shape[1] > 10:  # Only apply PCA if many features
                if hasattr(self.pca_transformer, "components_"):
                    features_scaled = self.pca_transformer.transform(features_scaled)
                else:
                    features_scaled = self.pca_transformer.fit_transform(
                        features_scaled
                    )

            return features_scaled

        except Exception as e:
            self.logger.error(f"Error preparing features: {e}")
            return None

    def _detect_statistical_anomalies(self, features_df: pd.DataFrame) -> List[int]:
        """Detect anomalies using statistical methods."""
        anomaly_indices = []

        try:
            numeric_columns = features_df.select_dtypes(include=[np.number]).columns

            for idx, row in features_df.iterrows():
                is_anomaly = False

                for col in numeric_columns:
                    value = row[col]
                    if pd.isna(value):
                        continue

                    # Calculate z-score for the column
                    col_mean = features_df[col].mean()
                    col_std = features_df[col].std()

                    if col_std > 0:
                        z_score = abs(value - col_mean) / col_std
                        if z_score > 3:  # 3-sigma rule
                            is_anomaly = True
                            break

                if is_anomaly:
                    anomaly_indices.append(idx)

            return anomaly_indices

        except Exception as e:
            self.logger.error(f"Error in statistical anomaly detection: {e}")
            return []

    def _combine_predictions(
        self,
        predictions: Dict[str, np.ndarray],
        scores: Dict[str, np.ndarray],
        features_df: pd.DataFrame,
        original_metrics: List[Dict],
        statistical_anomalies: List[int],
    ) -> List[AnomalyResult]:
        """Combine predictions from multiple algorithms using voting."""
        anomalies = []

        try:
            n_samples = len(features_df)

            for i in range(n_samples):
                votes = 0
                total_algorithms = 0
                confidence_scores = []

                # Count votes from ML algorithms
                for alg_name, pred in predictions.items():
                    if pred is not None and i < len(pred):
                        if pred[i] == -1:  # Anomaly
                            votes += 1
                        total_algorithms += 1

                        # Add confidence score
                        if alg_name in scores and scores[alg_name] is not None:
                            if i < len(scores[alg_name]):
                                confidence_scores.append(abs(scores[alg_name][i]))

                # Add statistical anomaly vote
                if i in statistical_anomalies:
                    votes += 1
                    confidence_scores.append(1.0)
                total_algorithms += 1

                # Decision based on majority vote
                if total_algorithms > 0:
                    vote_ratio = votes / total_algorithms

                    if vote_ratio >= 0.5:  # Majority vote for anomaly
                        # Calculate overall confidence
                        confidence = (
                            np.mean(confidence_scores) if confidence_scores else 0.5
                        )

                        # Create anomaly result
                        if i < len(original_metrics):
                            anomaly = self._create_anomaly_result(
                                original_metrics[i],
                                features_df.iloc[i],
                                confidence,
                                vote_ratio,
                                predictions,
                            )
                            if anomaly:
                                anomalies.append(anomaly)

            return anomalies

        except Exception as e:
            self.logger.error(f"Error combining predictions: {e}")
            return []

    def _create_anomaly_result(
        self,
        original_metric: Dict,
        feature_row: pd.Series,
        confidence: float,
        vote_ratio: float,
        algorithm_predictions: Dict[str, np.ndarray],
    ) -> Optional[AnomalyResult]:
        """Create an anomaly result object."""
        try:
            # Determine anomaly type based on metrics
            anomaly_type = self._classify_anomaly_type(original_metric, feature_row)

            # Determine severity based on confidence and deviation
            severity = self._calculate_severity(confidence, vote_ratio)

            # Extract relevant features
            features = {
                col: float(val) if pd.notna(val) else 0.0
                for col, val in feature_row.items()
                if pd.api.types.is_numeric_dtype(type(val))
            }

            return AnomalyResult(
                id=f"anomaly_{datetime.now().strftime('%Y%m%d_%H%M%S')}_{hash(str(original_metric)) % 10000}",
                type=anomaly_type,
                severity=severity,
                description=self._generate_anomaly_description(
                    original_metric, anomaly_type
                ),
                service=original_metric.get("labels", {}).get("service", "unknown"),
                namespace=original_metric.get("labels", {}).get("namespace", "default"),
                score=confidence,
                threshold=self.threshold,
                features=features,
                labels=original_metric.get("labels", {}),
                detected_at=datetime.now(),
                metadata={
                    "vote_ratio": vote_ratio,
                    "algorithms": list(algorithm_predictions.keys()),
                    "original_metric": original_metric,
                },
            )

        except Exception as e:
            self.logger.error(f"Error creating anomaly result: {e}")
            return None

    def _classify_anomaly_type(self, metric: Dict, features: pd.Series) -> str:
        """Classify the type of anomaly based on metric characteristics."""
        metric_name = metric.get("name", "").lower()

        if "cpu" in metric_name:
            return "resource_usage"
        elif "memory" in metric_name:
            return "resource_usage"
        elif "network" in metric_name:
            return "traffic_pattern"
        elif "error" in metric_name or "failed" in metric_name:
            return "error_rate"
        elif "latency" in metric_name or "duration" in metric_name:
            return "latency"
        elif "security" in metric_name or "auth" in metric_name:
            return "security"
        else:
            return "performance"

    def _calculate_severity(self, confidence: float, vote_ratio: float) -> str:
        """Calculate anomaly severity based on confidence and vote ratio."""
        combined_score = (confidence + vote_ratio) / 2

        if combined_score >= 0.9:
            return "critical"
        elif combined_score >= 0.7:
            return "high"
        elif combined_score >= 0.5:
            return "medium"
        else:
            return "low"

    def _generate_anomaly_description(self, metric: Dict, anomaly_type: str) -> str:
        """Generate a human-readable description of the anomaly."""
        metric_name = metric.get("name", "unknown")
        value = metric.get("value", 0)
        service = metric.get("labels", {}).get("service", "unknown")

        descriptions = {
            "resource_usage": f"Unusual resource usage detected in {metric_name} for service {service} (value: {value})",
            "traffic_pattern": f"Abnormal traffic pattern detected in {metric_name} for service {service}",
            "error_rate": f"Elevated error rate detected in {metric_name} for service {service}",
            "latency": f"Unusual latency pattern detected in {metric_name} for service {service}",
            "security": f"Security anomaly detected in {metric_name} for service {service}",
            "performance": f"Performance anomaly detected in {metric_name} for service {service}",
        }

        return descriptions.get(
            anomaly_type, f"Anomaly detected in {metric_name} for service {service}"
        )

    async def _train_models(self, features_df: pd.DataFrame):
        """Train all anomaly detection models."""
        try:
            self.logger.info("Training anomaly detection models...")

            features = self._prepare_features(features_df)
            if features is None:
                return

            # Train Isolation Forest
            if self.algorithms["isolation_forest"] is not None:
                self.algorithms["isolation_forest"].fit(features)

            # Train LOF
            if self.algorithms["lof"] is not None:
                self.algorithms["lof"].fit(features)

            # Initialize and train AutoEncoder
            if features.shape[1] > 0:
                self.algorithms["autoencoder"] = AutoEncoder(
                    hidden_neurons=[
                        features.shape[1] // 2,
                        features.shape[1] // 4,
                        features.shape[1] // 2,
                    ],
                    contamination=self.contamination,
                    epochs=50,
                    verbose=0,
                )
                self.algorithms["autoencoder"].fit(features)

            self.is_trained = True
            self.last_training = datetime.now()

            self.logger.info("Anomaly detection models trained successfully")

        except Exception as e:
            self.logger.error(f"Error training models: {e}")

    async def retrain_models(self):
        """Retrain models with recent data."""
        try:
            # Get recent training data
            recent_data = await self.data_processor.get_recent_training_data()

            if recent_data is not None and not recent_data.empty:
                await self._train_models(recent_data)
                self.logger.info("Models retrained successfully")
            else:
                self.logger.warning("No recent data available for retraining")

        except Exception as e:
            self.logger.error(f"Error retraining models: {e}")

    def is_ready(self) -> bool:
        """Check if the anomaly detector is ready to process requests."""
        return self.is_trained

    def get_health_status(self) -> Dict[str, Any]:
        """Get health status of the anomaly detector."""
        return {
            "healthy": self.is_trained,
            "last_training": (
                self.last_training.isoformat() if self.last_training else None
            ),
            "algorithms_loaded": len(
                [alg for alg in self.algorithms.values() if alg is not None]
            ),
            "total_algorithms": len(self.algorithms),
        }
