#!/usr/local/bin/jq -rf

.[]
# The LEGO parts of all the starter pack are 71200.
# Each 71170-71174 lists alternative building instructions which are from 71200.
| select(.productId > "71174")

| .productId as $productId
| .productName as $productName
| .productImage as $productImage

| .buildingInstructions[]
| select(.isAlternative)

| {
    productId: $productId,
    productName: $productName,
    productImage: $productImage,
    productImageWide: .frontpageInfo,
    model: .description,
    pdfLocation,
    # 71286 has filenames containing spaces
    file: .pdfLocation[(.pdfLocation | rindex("/"))+1:] | gsub("%20"; "_")
}