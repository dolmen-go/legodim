#!/usr/local/bin/jq -rf

# Usage:
#   ./list-models.jq data/index.json | ./fetch-files.jq | sh

.[]
| ("( "+
    "d='data/" + .productId + " " + (.productName | gsub("'"; "'\\''") ) + "'; " +
    "p=\"$d/" + .productId +"\"; " + 
    "f=\"$d\"/'" + (.model | gsub("'"; "'\\''")) + "'; " +
    "echo \"$f\"; " +
    "mkdir -p \"$d\"; " + 
    "{ test -f \"$p.png\" || curl -o \"$p.png\" '" + .productImageWide + "'; }; " +
    "{ test -f \"$p.jpg\" || curl -o \"$p.jpg\" '" + .productImage + "'; }; " +
    "{ test -f \"$f.pdf\" || curl -o \"$f.pdf\" '" + .pdfLocation + "'; }; " +
    "sips -i \"$p.png\" >/dev/null; " +
    "sips -i \"$p.jpg\" >/dev/null; " +
    "sips -i \"$f.pdf\" >/dev/null; " +
    ")")