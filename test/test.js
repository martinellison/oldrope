var gd = {};
var ld = {};
var currentPage = 'start';
var cp;
var pages;
var setPage = function(pageName) {
    console.log('displaying page: ' + pageName);
    currentPage = pageName;
    cp = pages[currentPage];
    if (!cp) console.error('unknown page: ' + currentPage);
    ld = cp.local;
    cp.set();
    console.log('displayed ' + currentPage);
};
var displayPage = function() {
    cp.redisplay();
    console.log('redisplayed ' + currentPage);
}
pages = {
    start: {
        local: {
            fred: false
        },
        set: function() {
            parts = [];
            parts.push('You are somewhere');
            parts.push('<a id=one>Go somewhere</a>');
            parts.push('<a id=two>Do something</a>');
            $('#main').html(parts.join("\n"));
            $('#one').click(function() {
                setPage('another');
                displayPage();
            });
            $('#two').click(function() {
                ld.fred = true;
                displayPage();
            });
        },
        redisplay: function() {
            if (ld.fred) $('#one').html('two');
        },
    },
    another: {
        local: {},
        set: function() {
            $('#main').html('end');
        },
        redisplay: function() {}
    },
};
setPage('start');
displayPage();
console.log('script loaded');