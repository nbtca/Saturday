const express = require("express");
const formidableMiddleware = require("express-formidable");
const router = express.Router();
const {  auth } = require("../middleware/auth");
const Element = require("../controller/element/element");

router.post("/login", Element.login);

router.use(auth);
router.get("/", Element.getAll);
router.get("/:rid", Element.get);
router.post("/", Element.create);
router.put("/", Element.update);
router.delete("/", Element.delete);
router.post("/update", Element.update);
router.post("/activate", Element.activate);
router.post("/updateAvatar", formidableMiddleware(), Element.updateAvatar);

module.exports = router;
