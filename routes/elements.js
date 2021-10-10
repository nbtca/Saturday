var express = require("express");
const router = express.Router();
const formidableMiddleware = require("express-formidable");
var Element = require("../controller/element/element");

router.get("/", Element.getAll);
router.get("/:rid", Element.get);
router.post("/", Element.create);
router.put("/", Element.update);
router.delete("/", Element.delete);
router.post("/update", Element.update);
router.post("/activate", formidableMiddleware(), Element.activate);
router.post("/updateAvatar", formidableMiddleware(), Element.updateAvatar);

module.exports = router;
