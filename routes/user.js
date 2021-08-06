var express = require('express');
var router = express.Router();
var host = 'rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com'
const mysql = require('serverless-mysql')({
    config: {
        host: host,
        user: 'user',
        password: '3ToMl5R2A7Fh',
        database: 'repairteam',
    }
})

/* GET home page. */
router.get('/',async function (req, res, next) {
    results = await mysql.query('SELECT count(*) FROM user')
    await mysql.end()
    res.send(results)
});

module.exports = router;
