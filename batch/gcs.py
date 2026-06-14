from pathlib import Path
from google.cloud import storage
from google.oauth2 import service_account
import os
import pandas as pd
from io import BytesIO

def check_credentials() -> bool:
    # GCP上ではクレデンシャルはいらないので、ローカルでの実行時のみクレデンシャルの存在を確認する
    check = os.getenv("GCS_SERVICE_ACCOUNT_KEY_PATH")
    if check is None:
        print("[INFO] GCS_SERVICE_ACCOUNT_KEY_PATH is not set.")
        return False
    if not Path(check).exists():
        print(f"[ERROR] GCS service account key file not found at: {check}")
        return False
    return True

def fetch_gcs_file_to_df(gcp_path: str) -> pd.DataFrame:
    if check_credentials():
        credentials = service_account.Credentials.from_service_account_file(os.getenv("GCS_SERVICE_ACCOUNT_KEY_PATH"))
        client = storage.Client(credentials=credentials)
    else:
        # GCP上で実行されている場合はクレデンシャルが不要
        client = storage.Client()
        print("[WARNING] GCS credentials not found. Skipping file fetch.")

    bucket_name = os.getenv("GCS_BUCKET_NAME")
    if bucket_name is None:
        raise ValueError("GCS_BUCKET_NAME environment variable is not set.")
    
    bucket = client.bucket(bucket_name)
    blob = bucket.blob(gcp_path)

    df = pd.read_csv(BytesIO(blob.download_as_string()))

    return df


def upload_df_to_gcs(df: pd.DataFrame, gcp_path: str) -> None:
    if check_credentials():
        credentials = service_account.Credentials.from_service_account_file(os.getenv("GCS_SERVICE_ACCOUNT_KEY_PATH"))
        client = storage.Client(credentials=credentials)
    else:
        # GCP上で実行されている場合はクレデンシャルが不要
        client = storage.Client()
        print("[WARNING] GCS credentials not found.")

    bucket_name = os.getenv("GCS_BUCKET_NAME")
    if bucket_name is None:
        raise ValueError("GCS_BUCKET_NAME environment variable is not set.")
    
    bucket = client.bucket(bucket_name)
    blob = bucket.blob(gcp_path)

    blob.upload_from_string(df.to_csv(index=False), content_type='text/csv')