#!/usr/bin/env node
const WebSocket = require('ws');
const jwt = require('jsonwebtoken');
const spawn = require('child_process').spawn;
const secret = "hellloafsdfefss9843ru93f";

const run = (channel1) => {
    let address = process.env.ADDRESS || "ws://ws.z.12zan.net/dbcast/some_channel_as_you_wish";
    try {
        let ws = new WebSocket(address);
        ws.on('open', function open() {
            //setInterval(function () {
                //ws.send(`一条消息` + channel1);
                //之后就发一条消息；只有这个组内的所有客户端能收到;
            //}, 1500);
        });
        ws.on('message', function incoming(data) {
            console.log("incoming data:");
            sign = jwt.decode(data,secret);
            console.log(JSON.parse(sign.sub));
            //收到消息，打印出来;
        });
        ws.on("error", function (e) {
            console.log("got error");
            console.log(e);
        });
    } catch (e) {
        console.log("got new exception of client.");
        console.log(e);
    }
};
run();
