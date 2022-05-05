const express = require("express");
const path = require("path");
const cookieParser = require("cookie-parser");
const logger = require("morgan");
const app = express();

app.use(logger("dev"));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, "public")));

const main = require("./routes/main");
app.use(
  "/",
  (req, res, next) => {
    console.log(req.body);
    next();
  },
  main
);

// catch 404 and forward to error handler
app.use(function (req, res, next) {
  res.status(404).send('Sorry cant find that!');
});

// error handler
app.use(function (err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get("env") === "development" ? err : {};
  res.status(500).send('Something broke!');
});

module.exports = app;
