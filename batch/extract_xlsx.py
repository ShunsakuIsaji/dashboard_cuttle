import pandas as pd
import re

def infer_year_from_month(month: int, periods: int = 15) -> int:
    """
    月だけが与えられた場合、現在の月と比較して開始年を推測する
    例えば、現在が2026年6月で、monthが1月なら、2025年1月と推測する
    """
    now = pd.Timestamp.now()
    current_year = now.year
    
    candidate_year = current_year
    while True:
        if pd.Timestamp(year=candidate_year, month=month, day=1) + pd.DateOffset(months=periods-1) > now:
            candidate_year -= 1
        else:
            break
    
    return candidate_year

def parse_start_date(value) -> pd.Timestamp:
    """
    '2024年1月' のような値から月初日 Timestamp を作る
    1. エクセルのシリアル値パターン
    2. 'YYYY年M月' パターン
    3. 'M月' パターン（年は infer_year_from_monthで推測）
    """
    if isinstance(value, (int, float)):
        # エクセルのシリアル値を日付に変換
        if value < 13:
            # 月だけの可能性があるので、年は推測する
            month = int(value)
            year = infer_year_from_month(month)
            return pd.Timestamp(year=year, month=month, day=1)
        return pd.Timestamp("1899-12-30") + pd.to_timedelta(value, unit='D')

    text = str(value).strip()
    print(f"[DEBUG] Parsing start date from text: {text}")
    m = re.search(r"(\d{4})年\s*(\d{1,2})", text)
    if not m:
        k = re.search(r"(\d{1,2})", text)
        if k:
            # Handle case where only month is specified
            month = int(k.group(1))
            year = infer_year_from_month(month)
            return pd.Timestamp(year=year, month=month, day=1)
        raise ValueError(f"start date is not YYYY年M月 format: {value}")

    year = int(m.group(1))
    month = int(m.group(2))
    print(f"[DEBUG] Extracted year: {year}, month: {month}")
    return pd.Timestamp(year=year, month=month, day=1)


def parse_dates(date_series: pd.Series, periods: int = 15) -> list[str]:
    """
    先頭行の 'yyyy年n月' を起点に、periodsか月分の YYYY-MM を返す
    """
    non_empty = date_series.dropna()
    if non_empty.empty:
        raise ValueError("date column is empty")

    start = parse_start_date(non_empty.iloc[0])

    return [
        (start + pd.DateOffset(months=i)).strftime("%Y-%m")
        for i in range(periods)
    ]

def extract_data_from_xlsx(file_path: str, sheet_name: str, skip_rows: int, date_column: int, items: dict[str, dict], periods: int = 15) ->  pd.DataFrame:
    df = pd.read_excel(file_path, sheet_name=sheet_name, skiprows=skip_rows, header=None)
    dates = parse_dates(df.iloc[:, date_column], periods)

    records = []
    for item_key, item_info in items.items():
        category = item_info.get("category", "uncategorized")
        data_column = item_info["data_column"]
        unit = item_info.get("unit", "")

        for date, value in zip(dates, df.iloc[:periods, data_column]):
            if pd.isna(value):
                continue
            records.append({
                "date": date,
                "category": category,
                "item": item_key,
                "value": value,
                "unit": unit,
            })
    
    return pd.DataFrame(records)