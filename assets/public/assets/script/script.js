window.addEventListener("load", function () {
    document.addEventListener("keyup", function (e) {
        if (e.ctrlKey && e.key === "r") {
            window.location.reload();
        }
    });
});

window.addEventListener("resize", function() {
    clearTimeout(resizeHandle);

    resizeHandle = setTimeout(onResize, 750);
});

var resizeHandle;

function onResize() {
    post("/api/window/save", JSON.stringify({
        "window_width": window.innerWidth,
        "window_height": window.innerHeight,
    }));
}

function path() {
    var url = window.location.href;
    var start = url.indexOf("#");
    if (start === -1) {
        return "";
    }
    var end = url.indexOf("?", start);
    if (end === -1) {
        end = url.length;
    }
    return url.substring(start + 1, end);
}

function nav(url) {
    window.location.href = uniqify(url);
}

function get(url, callback) {
    request("GET", url, null, callback);
}

function post(url, body, callback) {
    request("POST", url, body, callback);
}

function del(url, callback) {
    request("DELETE", url, null, callback)
}

function uniqify(url) {
    if (url.lastIndexOf("?") === -1) {
        return url + "?ts=" + Date.now();
    } else {
        return url + "&ts=" + Date.now();
    }
}

function request(method, url, body, callback) {
    var req = new XMLHttpRequest();
    req.open(method, uniqify(escape(url)), true);
    req.setRequestHeader("Cache-Control", "no-cache, must-revalidate, post-check=0, pre-check=0");
    req.setRequestHeader("Cache-Control", "max-age=0");
    req.setRequestHeader("Pragma", "no-cache");
    req.setRequestHeader("Vary", "*");
    req.onload = function () {
        if (this.status !== 200) {
            return
        }
        if (this.readyState !== 4) {
            return;
        }
        if (callback) {
            callback(JSON.parse(req.responseText));
        }
    };
    req.send(body);
}