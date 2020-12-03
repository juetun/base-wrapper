var pageFuncObj = {

    changeBrand: function (ob, bHid) {
        var o = $("#navbarDropdown1");
        o.html(ob.attr("bname")).attr({select_id: bHid});
        $.ajax({
            url: o.attr("load_url"), type: "GET", data: {brand_hid: bHid}, async: true, success: function (res) {
                var s = "";
                for (var i in res.data.sty_list) {
                    var dt = res.data.sty_list[i];
                    s += "<a class=\"dropdown-item\" onclick='pageFuncObj.changeStyle($(this),\"" + dt.hid + "\")' href=\"javascript:;\">" + dt.zh_name + "</a>";
                }
                $("#navbarDropdownmenu2").html(s);
            }
        });
    },
    changeStyle: function (ob, hid) {
        var o = $("#navbarDropdown2");
        o.html(ob.html()).attr({select_id: hid});
        $.ajax({
            url: o.attr("load_url"), type: "GET", data: {style_hid: hid}, async: true, success: function (res) {
                console.log(res);
                var s = "";
                for (var i in res.data.model_list) {
                    var dtm = res.data.model_list[i];
                    s += "<a class=\"dropdown-item  bg-secondary text-light pt-2 pb-2 pl-2\" href=\"javascript:;\">" + dtm.model_label.label + "</a>";
                    for (var j in dtm.model_list) {
                        var dt = dtm.model_list[j];
                        s += "<a class=\"dropdown-item\" onclick='pageFuncObj.changeModel($(this),\"" + dt.hid + "\")' href=\"javascript:;\">" + dt.whole_name + "</a>";
                    }
                }
                $("#navbarDropdownmenu3").html(s);
            }
        });
    },
    changeModel: function (ob, hid) {
        var o = $("#navbarDropdown3");
        o.html(ob.html()).attr({select_id: hid});
    },
};

function searchCar(o) {
    var o3 = $("#navbarDropdown3");
    if (o3.attr("select_id")) {
        location.href = o3.attr("location").replace("hid.html", o3.attr("select_id") + '.html');
        return;
    }
    var o2 = $("#navbarDropdown2");
    if (o2.attr("select_id")) {
        location.href = o2.attr("location").replace("hid.html", o2.attr("select_id") + '.html');
        return;
    }
    var o1 = $("#navbarDropdown1");
    if (o1.attr("select_id")) {
        location.href = o1.attr("location").replace("hid.html", o1.attr("select_id") + '.html');
        return;
    }
}

$(window).scroll(function () {
    //为了保证兼容性，这里取两个值，哪个有值取哪一个
    var ob = $("#hot-article-bottom");

    //scrollTop就是触发滚轮事件时滚轮的高度
    var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
    var obt = ob.offset().top - $(window).height();
    if (obt > $(window).height() && scrollTop > obt) {
        setTimeout(function () {
            if (!parseInt(ob.attr("loading")) && !parseInt(ob.attr("finish"))) {
                ob.attr({"loading": 1}).find(".child").removeClass("d-none");
                var page = parseInt(ob.attr("page_no")) + 1;
                ob.attr("page_no", page);
                if (page > 4) { //最多加载50页数据
                    ob.attr({"loading": 0, "finish": 1}).find(".child").addClass("d-none");
                    return;
                }
                getArticle(ob.attr("load_url"), {
                    page_no: page,
                    page_size: parseInt(ob.attr("page_size"))
                }, function (res) {
                    if (res.trim() === "") {
                        ob.attr({"loading": 0, "finish": 1}).find(".child").addClass("d-none");
                        return
                    }
                    $("#article-list-load").append(res);
                    ob.attr({"loading": 0, finish: 0}).find(".child").addClass("d-none");
                });
            }
        },500);

    }
    // console.log("滚动距离" + scrollTop,$("#hot-article-bottom").offset().top-$(window).height());//+$(window).height()
});

function getArticle(url, arg, callBack) {
    $.ajax({
        url: url, type: "GET", data: arg, async: true, success: function (res) {
            callBack(res);
        }
    });
}