var express = require("express");
const router = express.Router();
const Element = require("../controller/element/element");

router.get("/", Element.getAll);
router.get("/:rid", Element.get);
router.post("/", Element.create);
router.put("/", Element.update);
router.delete("/", Element.delete);

module.exports = router;
