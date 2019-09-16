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
    post("/api/save/window", JSON.stringify({
        "window_width": window.innerWidth,
        "window_height": window.innerHeight,
    }));
}

function path() {
    return window.location.href.split("#")[1] || "";
}

function nav(url) {
    window.location.href = url;
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

function request(method, url, body, callback) {
    var req = new XMLHttpRequest();
    req.open(method, url, true);
    req.setRequestHeader('cache-control', 'no-cache, must-revalidate, post-check=0, pre-check=0');
    req.setRequestHeader('cache-control', 'max-age=0');
    req.setRequestHeader('pragma', 'no-cache');
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