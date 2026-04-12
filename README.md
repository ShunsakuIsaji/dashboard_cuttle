# dashboard_cuttle


## ページ構成

- `/index`
- `/calf` - 子牛
- `/milk` - 酪農生乳
- `/fattening` - 肥育
- `/feed` - 飼料ほか農資材
- `/indices` - 他指標（原油価格、穀物先物、バルチック海運）

## データソース

- [農業物価統計調査](https://www.e-stat.go.jp/stat-search/files?page=1&layout=datalist&toukei=00500204&tstat=000001170266&cycle=7&tclass1=000001201980&tclass2=000001230427&tclass3val=0&metadata=1&data=1)
    - 農産物価格（年次:2020-2024）
    - 農資材価格（年次:2020-2024）

- [農畜産業振興機構](https://www.alic.go.jp/livestock/index.html)
    - [統計資料：牛乳・乳製品](https://www.alic.go.jp/joho-c/raku02_000085.html)
        - [平均泌乳量](https://lin.alic.go.jp/alic/statis/dome/data2/j_excel/5030a.xlsx)
        - [生乳農家販売価格](https://lin.alic.go.jp/alic/statis/dome/data2/j_excel/5051A.xlsx)
    - [統計資料：肥育・飼料](https://www.alic.go.jp/joho-c/joho05_000073.html)
        - [枝肉の規格別卸売価格](https://lin.alic.go.jp/alic/statis/dome/data2/j_excel/2060a.xlsx)
        - [子牛の取引頭数と価格](https://lin.alic.go.jp/alic/statis/dome/data2/j_excel/2140a.xlsx)
        - [飼料の輸入価格](https://lin.alic.go.jp/alic/statis/dome/data2/j_excel/7011b.xlsx)
        - [配合飼料の価格](https://lin.alic.go.jp/alic/statis/dome/data2/j_excel/7022a.xlsx)

|category|item|unit|説明|ソース|
|:-------|:---|:---|:---|:----|
|milk|milk_yield_hokkaido|kg/日/頭|北海道の1日あたり平均泌乳量|平均泌乳量|
|milk|milk_yield_honshu|kg/日/頭|本州の1日あたり平均泌乳量|平均泌乳量|
|milk|milk_price|円/kg|生乳農家販売価格|生乳農家販売価格|
|fattening|wagyu_cow_a5|円/kg|和牛メスA5等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|wagyu_cow_a4|円/kg|和牛メスA4等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|wagyu_cow_a3|円/kg|和牛メスA3等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|wagyu_steer_a5|円/kg|和牛去勢A5等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|wagyu_steer_a4|円/kg|和牛去勢A4等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|wagyu_steer_a3|円/kg|和牛去勢A3等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|wagyu_steer_a2|円/kg|和牛去勢A2等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|dairy_cow_c2|円/kg|乳牛メスC2等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|cross_cow_b3|円/kg|交雑メスB3等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|cross_cow_b2|円/kg|交雑メスB2等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|cross_steer_b3|円/kg|交雑去勢B3等級の枝肉価格|枝肉の規格別卸売価格|
|fattening|cross_steer_b2|円/kg|交雑去勢B2等級の枝肉価格|枝肉の規格別卸売価格|
|calf|wagyu_female|千円/頭|黒毛和種メスの平均価格|子牛の取引頭数と価格|
|calf|wagyu_male|千円/頭|黒毛和種オスの平均価格|子牛の取引頭数と価格|
|calf|holstein_female|千円/頭|ホルスタイン種メスの平均価格|子牛の取引頭数と価格|
|calf|holstein_male|千円/頭|ホルスタイン種オスの平均価格|子牛の取引頭数と価格|
|calf|cross_female|千円/頭|交雑種メスの平均価格|子牛の取引頭数と価格|
|calf|cross_male|千円/頭|交雑種オスの平均価格|子牛の取引頭数と価格|
|feed|corn_price|円/t|トウモロコシの輸入価格|飼料の輸入価格|
|feed|sorghum_price|円/t|こうりゃんの輸入価格|飼料の輸入価格|
|feed|barley_price|円/t|大麦の輸入価格|飼料の輸入価格|
|feed|wheat_price|円/t|小麦の輸入価格|飼料の輸入価格|
|feed|soy_price|円/t|大豆油かすの輸入価格|飼料の輸入価格|
|feed|fishmeal_price|円/t|魚粉の輸入価格|飼料の輸入価格|
|feed|hay_price|円/t|乾牧草の輸入価格|飼料の輸入価格|
|feed|haycube_price|円/t|ヘイキューブの輸入価格|飼料の輸入価格|
|feed|compound_price|円/t|配合飼料の工場渡価格|配合飼料の価格|


'''mermaid
erDiagram
    prices{
        date date
        string source
        string category
        string item
        int value
        string unit
    }   
'''
