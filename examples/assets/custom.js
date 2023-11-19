function myFunction() {
    var txt;

    if (confirm("Please Press a button!")) {
        txt = "You pressed Button!";
    } else {
        txt = "You pressed Cancel Button!";
    }

    document.getElementById("Test Confirm Box").innerHTML = txt;
}

