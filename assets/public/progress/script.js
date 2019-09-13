function status(text) {
    document.getElementById("status").innerText = text;
}

function message(text) {
    document.getElementById("message").innerText = text;
}

function overall(progress) {
    if (progress >= 0) {
        document.getElementById("overall").setAttribute("value", progress);
    } else {
        document.getElementById("overall").removeAttribute("value");
    }
}

function task(progress) {
    document.getElementById("task").classList.remove("hidden");
    if (progress >= 0) {
        document.getElementById("task").setAttribute("value", progress);
    } else {
        document.getElementById("task").removeAttribute("value");
    }
}