# Go interface to LEGO Dimensions Toy Pad

## Status

Work in progress.

Example: see cmd/toypad-read/main.go

## Resources

### Toy Pad
* https://github.com/Lefinnois/legopad_hid/blob/master/README.protocol.md
* https://github.com/nathankellenicki/node-legodimensions/blob/master/src/toypad.js
* https://github.com/ags131/node-ld/ Node.js
* https://www.dajlab.org/jtoypad.html (Java code: https://github.com/e-amzallag/jtoypad)
* https://github.com/VincentXII/Lego-Dimensions-Pad-Scripts Python
* https://github.com/liloumuloup/Lego-Dimensions-Pad-Scripts/blob/master/legousb-partytime-random.py Python
* https://github.com/mpetrov/node-dimensions Node.js
* Communauté française: https://forum.freelug.org/viewtopic.php?f=10&t=216&sid=3140983d0fbe23757469ecb8b05a6919&start=15

### Tags
* https://nfc.toys/ (talks about Disney Infinity which is the same protocol as LEGO Dimensions)
* https://nfc.toys/interop-inf.html Code
* https://www.reddit.com/r/Legodimensions/comments/3oixfa/for_those_having_trouble_with_toy_or_figure_tags/ Reset tags
* http://www.proxmark.org/forum/viewtopic.php?id=2657
* https://nfc-bank.com/bins.php?categoryid=28
* https://pastebin.com/mB5zrtxx
* Vehicules ids: https://pastebin.com/KUBRtaxi
* List of characters & vehicles: https://github.com/drake-vcu/LegoCrypto/tree/master/LegoCrypto.Data/Resources
* https://github.com/phogar/ldnfctags/
    * https://github.com/phogar/ldnfctags/blob/master/src/legodimensions_characters.c
    * https://github.com/phogar/ldnfctags/blob/master/src/legodimensions_vehicles.c
* Vehicles/Acessories are "empty by default but are programmed by the LEGO Dimensions game". Source: [WB Games FAQ](https://legogamessupport.wbgames.com/hc/fr-fr/articles/360001217048-Comment-reprogrammer-votre-Toy-Tag-de-v%C3%A9hicule-ou-d-accessoire-LEGO-)
* Vehicles/Accesories are "rewritable RFID". Source: [WB Games FAQ](https://legogamessupport.wbgames.com/hc/fr-fr/articles/360001216968-Comment-fonctionnent-les-Toy-Tags-et-les-am%C3%A9liorations-de-gadgets-et-de-v%C3%A9hicules-)

### Building instructions
https://www.lego.com/fr-fr/service/buildinginstructions/search#?search&theme=10000-20229

### Game
* Packattack Guide on YouTube https://www.youtube.com/playlist?list=PLQEhT4w332jqqpysL0ngF1EOg9ragoFow
* Soluce en français https://www.supersoluce.com/soluce/lego-dimensions/soluce-lego-dimensions

### Communities
* [LEGO Dimensions on Reddit](https://www.reddit.com/r/Legodimensions/)

### Official support
* LEGO: https://www.lego.com/fr-fr/service/help/produits/themes-et-ensembles/dimensions
* WB Games: https://legogamessupport.wbgames.com/hc/fr-fr/categories/360000090167-LEGO-Dimensions

### LEGO Dimensions Collection Vortex
* https://lego-dimensions.fandom.com/wiki/LEGO_Dimensions_Collection_Vortex
* [Android](https://play.google.com/store/apps/details?id=com.wb.lego.dimensions)
* [Download Android APK](https://apkpure.com/lego%C2%AE-dimensions%E2%84%A2/com.wb.lego.dimensions)




### Description of Tag data
#### Page 38
* 00000000 => Character
* 00010000 => Vehicle

#### Page 45
Tag UID : 010203??
#### Page 46
Tag UID : 04050607