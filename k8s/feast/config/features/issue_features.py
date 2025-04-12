"""
Feature definitions for issue data.
"""

from datetime import timedelta
from feast import Entity, Feature, FeatureView, Field, FileSource
from feast.types import Float32, Int64, String

issue = Entity(
    name="issue",
    join_keys=["issue_id"],
    description="Issue identifier",
)

issue_source = FileSource(
    path="s3://feast/data/issue_features.parquet",
    timestamp_field="timestamp",
    created_timestamp_column="created_timestamp",
)

issue_features = FeatureView(
    name="issue_features",
    entities=[issue],
    ttl=timedelta(days=365),
    schema=[
        Field(name="issue_id", dtype=Int64),
        Field(name="repository", dtype=String),
        Field(name="title_embedding", dtype=Float32, shape=(384,)),
        Field(name="description_embedding", dtype=Float32, shape=(384,)),
        Field(name="topic_vector", dtype=Float32, shape=(50,)),
        Field(name="language", dtype=String),
        Field(name="stars", dtype=Int64),
        Field(name="issue_age_days", dtype=Int64),
        Field(name="solution_length", dtype=Int64),
        Field(name="has_code", dtype=Int64),
    ],
    source=issue_source,
    online=True,
    tags={"team": "ml", "owner": "data-science"},
)

repository_source = FileSource(
    path="s3://feast/data/repository_features.parquet",
    timestamp_field="timestamp",
    created_timestamp_column="created_timestamp",
)

repository_features = FeatureView(
    name="repository_features",
    entities=[issue],
    ttl=timedelta(days=365),
    schema=[
        Field(name="issue_id", dtype=Int64),
        Field(name="repository", dtype=String),
        Field(name="repository_embedding", dtype=Float32, shape=(384,)),
        Field(name="repository_stars", dtype=Int64),
        Field(name="repository_forks", dtype=Int64),
        Field(name="repository_age_days", dtype=Int64),
        Field(name="repository_topics", dtype=String, shape=(10,)),
    ],
    source=repository_source,
    online=True,
    tags={"team": "ml", "owner": "data-science"},
)
