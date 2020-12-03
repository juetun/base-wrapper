var commonObj = {
        compareKey: "compareModel",
        limitCount: 8,
        showCompare: function (u) {
            var dt = this.haveData();
            if (dt.length < 2) {
                $("#exampleModal").modal('show').find(".modal-body").html(`请至少选择两种车型`);
                return
            }
            var hidArray = [];
            for (var i in dt) {
                hidArray.push(dt[i].hid);
            }
            u += this.urlEncode(hidArray.join(","));
            this.clearData();
            location.href = u;
        },
        urlEncode: function (str) {
            str = (str + '').toString();
            return encodeURIComponent(str).replace(/!/g, '%21').replace(/'/g, '%27').replace(/\(/g, '%28').replace(/\)/g, '%29').replace(/\*/g, '%2A').replace(/%20/g, '+');
        },
        addCompare: function (hid, title) {
            var dt = this.haveData();
            $("#compare-count").removeClass("d-none");
            if (dt.length > this.limitCount) {
                $("#exampleModal").modal('show').find(".modal-body").html(`对比车型最多选择${this.limitCount}种车型`);
                return;
            }
            for (var i in dt) {
                if (dt[i].hid == hid) {
                    $("#compare-num,#compare-num-1").html(dt.length);
                    $("#exampleModal").modal('show').find(".modal-body").html(`您已添加过该车型`);
                    return
                }
            }
            dt.push({hid: hid, t: title});
            $.cookie(this.compareKey, JSON.stringify(dt), {expires: 7, path: '/'});
            this.initCompare();
        }
        ,
        clearData: function () {
            $.cookie(this.compareKey, null, {expires: 7, path: '/'});
        },
        haveData: function () {
            var c = $.cookie(this.compareKey);
            var dt = [];
            if (c) {
                dt = $.parseJSON(c);
                if (dt == null) {
                    dt = [];
                }
            }
            return dt;
        }
        ,
        delete: function (hid) {
            // console.log(hid);
            var tmp = this.haveData();
            var dt = [];
            for (var i in tmp) {
                if (tmp[i].hid == hid) {
                    continue;
                }
                dt.push(tmp[i]);
            }
            $.cookie(this.compareKey, JSON.stringify(dt), {expires: 7, path: '/'});
            this.initCompare();
        },
        getDropMenu: function (dt) {
            var s = "";
            for (var i in dt) {
                s += ` <li class="list-group-item  border-0"> \n
                            ${dt[i].t}&nbsp;
                            <span onclick=commonObj.delete('` + dt[i].hid + `');return; class="badge badge-primary badge-pill float-right delete-compare">✕</span>
                        </li>`;
            }
            return s
        },
        initCompare: function () {
            var dt = this.haveData();
            var ob = $("#compare-num,#compare-num-1");
            if (dt.length == 0) {
                ob.addClass("d-none");
            }
            ob.html(dt.length);
            if (dt.length > 0) {
                ob.removeClass("d-none");
                $("#car-brand-pinyin").removeClass("d-none");
                $("#compare-count").removeClass("d-none");
                $("#dropdownCompareDMenu").html(this.getDropMenu(dt));
            } else {
                $("#compare-count").addClass("d-none");
                $("#car-brand-pinyin").addClass("d-none");

            }
        },
    }
;
$(function () {
    commonObj.initCompare();
    $('[data-toggle="tooltip"]').tooltip();
});
