
function r () {
    return `hello ${function(){var res = "";for (var i = 0; i < 10; i++) {res+="fishes"};return res;}()}`
}

console.log(r())