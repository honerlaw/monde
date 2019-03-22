window.onload = () => {
    document.querySelector('form').addEventListener('submit', (e) => {
        document.querySelector("form button").innerHTML = "uploading...";
        document.querySelector("form button").setAttribute("disabled", true);
    });
}
