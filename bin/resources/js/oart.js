var g_curGroupName;
var g_action;
var g_nodelist = {};
var g_createTableData;
var g_selectedNode;
var g_spinner = null
var g_all_Templates = {};
var g_all_ConfigTemplate = {}
var g_all_ConfigTemplate_renderObj = {}
var g_all_ProcessInfo = {}
var g_resourcelist = {}
var g_resourceid = 1

var g_ServerId = 1
var g_Uid = 20001
var g_IpPort = "10.20.104.50:20003"
var g_AlarmPageNum = 0
var g_oldRightList = ""
var g_fileContent = {}
var g_fileName = ""
var g_levels = {"4": true, "0": true}
var g_PageNum = 1
var g_PageSize = 1000+""
function sleep(numberMillis) {
    var now = new Date();
    var exitTime = now.getTime() + numberMillis;
    while (true) {
        now = new Date();
        if (now.getTime() > exitTime)
            return;
    }
}

function onClickResourceManager() {
    myajax("recourcelist", "")
}

function onClickAddAlarm() {
    myajaxV2("recourcelist", "", "onclickaddalarm")
}

function onClickShowAlarm() {
    myajax("getAlarmList", "0")
}

function onGetAlarmByPageNum(nPageNum) {
    myajax("getAlarmList", nPageNum.toString())
    g_AlarmPageNum = nPageNum
}

function onGetAlarmPre() {
    if ((g_AlarmPageNum - 1) < 0) {
        g_AlarmPageNum = 0
    } else {
        g_AlarmPageNum = g_AlarmPageNum - 1
    }
    myajax("getAlarmList", g_AlarmPageNum.toString())
}

function onGetAlarmNext() {
    g_AlarmPageNum = g_AlarmPageNum + 1
    myajax("getAlarmList", g_AlarmPageNum.toString())
}

function Index2HourMinute(nIndex) {
    var Hour = Math.floor(nIndex * 5 / 60)
    var Minute = nIndex * 5 % 60
    return Hour.toString() + ":" + Minute.toString()
}

function Plant(DataSlices) {
//为data准备的数据
    var mylables = []
    for (i = 0; i < 24; i++) {
        mylables.push(i.toString())
    }
//用于存放图表上的数据
    var lineChartData = {
        labels: mylables,
        //数据集（y轴数据范围随数据集合中的data中的最大或最小数据而动态改变的）
        datasets: [
            {
                fillColor: "transparent",     //背景色，常用transparent透明
                strokeColor: "rgba(220,0,0,1)",  //线条颜色，也可用"#ffffff"
                pointColor: "rgba(220,0,0,1)",   //点的填充颜色
                pointStrokeColor: "#fff",            //点的外边框颜色
                lable: "当天",
                data: []      //点的Y轴值
            },

            {
                fillColor: "transparent",
                strokeColor: "rgba(0,0,205,1)",
                pointColor: "rgba(0,0,205,1)",
                pointStrokeColor: "#fff",
                lable: "前一天",
                data: []    //data中的参数，通过下方for循环，获取arr2中的数据
            }
        ]
    }

    lineChartData.datasets[0].data = DataSlices[0];
    lineChartData.datasets[1].data = DataSlices[1];

//定义图表的参数
    var defaults = {
//网格线是否在数据线的上面
        scaleOverlay: false,

        //是否用硬编码重写y轴网格线
        scaleOverride: false,

        //** Required if scaleOverride is true **
        //y轴刻度的个数
        scaleSteps: null,

        //y轴每个刻度的宽度
        scaleStepWidth: 20,

        // Y 轴的起始值
        scaleStartValue: null,
        // Y/X轴的颜色
        scaleLineColor: "rgba(0,0,0,.1)",
        // X,Y轴的宽度
        scaleLineWidth: 1,
        // 刻度是否显示标签, 即Y轴上是否显示文字
        scaleShowLabels: false,
        // Y轴上的刻度,即文字
        scaleLabel: "<%=value%>",
        // 字体
        scaleFontFamily: "'Arial'",
        // 文字大小
        scaleFontSize: 16,
        // 文字样式
        scaleFontStyle: "normal",
        // 文字颜色
        scaleFontColor: "#666",
        // 是否显示网格
        scaleShowGridLines: true,
        // 网格颜色
        scaleGridLineColor: "rgba(0,0,0,.05)",
        // 网格宽度
        scaleGridLineWidth: 2,
        // 是否使用贝塞尔曲线? 即:线条是否弯曲
        bezierCurve: true,
        // 是否显示点数
        pointDot: false,
        // 圆点的大小
        pointDotRadius: 2,
        // 圆点的笔触宽度, 即:圆点外层白色大小
        pointDotStrokeWidth: 2,
        // 数据集行程(连线路径)
        datasetStroke: true,
        // 线条的宽度, 即:数据集
        datasetStrokeWidth: 2,
        // 是否填充数据集
        datasetFill: true,
        // 是否执行动画
        animation: true,
        // 动画的时间
        animationSteps: 60,
        // 动画的特效
        animationEasing: "easeOutQuart",
        // 动画完成时的执行函数
        /*onAnimationComplete: null*/

    }
    Chart.defaults.global.pointHitDetectionRadius = 3
    var ctx = document.getElementById("myChart").getContext("2d");
    window.myLine = new Chart(ctx).Line(lineChartData, defaults);
}

function OnClickSaveResource(id) {
    var requestTable = {}
    requestTable.Type = parseInt($("#Type_" + id).val())
    requestTable.Id = parseInt($("#Id_" + id).val())
    requestTable.Name = $("#Name_" + id).val()
    requestTable.FiveMinuteThreshold = parseInt($("#FiveMinuteThreshold_" + id).val())
    requestTable.OneHourThreshold = parseInt($("#OneHourThreshold_" + id).val())
    requestTable.OneDayThreshold = parseInt($("#OneDayThreshold_" + id).val())
    requestTable.Filter = $("#Filter_" + id).val()
    myajax("saveResource", JSON.stringify(requestTable))
    console.log(JSON.stringify(requestTable))
}

function OnClickShowResourceSlices() {
    /*$('#temp_template').html(g_all_Templates["ResourceSlicesQuery.html"])
    var source = $('#ResourceSlice-template').html();
    var myTemplate = Handlebars.compile(source);
    var out = myTemplate({})
    $('#rightlist').html(out);*/
    myajaxV2("recourcelist", "", "OnClickShowResourceSlices")
    //Plant([])
    //Datetime()
    //$(".chzn-select").chosen(); $(".chzn-select-deselect").chosen({allow_single_deselect:true});
}

function Datetime() {
    $('#datetimepicker1').datetimepicker({
        language: 'zh-CN',//显示中文
        format: 'yyyy-mm-dd',//显示格式
        minView: "month",//设置只显示到月份
        initialDate: new Date(),
        autoclose: true,//选中自动关闭
        todayBtn: true,//显示今日按钮
        locale: moment.locale('zh-cn')
    });
    //默认获取当前日期
    var today = new Date();
    var nowdate = (today.getFullYear()) + "-" + (today.getMonth() + 1) + "-" + today.getDate();
    //对日期格式进行处理
    var date = new Date(nowdate);
    var mon = date.getMonth() + 1;
    var day = date.getDate();
    var mydate = date.getFullYear() + "-" + (mon < 10 ? "0" + mon : mon) + "-" + (day < 10 ? "0" + day : day);
    document.getElementById("nowdate").value = mydate;
}

function OnQueryResourceSlices() {
    var strDay = document.getElementById("nowdate").value
    strDay = strDay.replace(/-/g, "")
    var strResourceId = $('#example28').val()

    myajax("getResourceSlice", strDay + "|" + strResourceId)
    //Plant()
}

function myajax(action, data) {
    g_action = action;
    strUrl = "ajax?" + "action=" + action + "&data=" + data;
    //$.history.load(url,dealWithAjaxRetData)
    $.get(strUrl, dealWithAjaxRetData);
}

function myajaxV2(action, data, afterActionName) {
    g_action = afterActionName;
    strUrl = "ajax?" + "action=" + action + "&data=" + data;
    $.get(strUrl, dealWithAjaxRetData);
}

function UnicodeToUtf8(unicode) {
    var uchar;
    var utf8str = "";
    var i;

    for (i = 0; i < unicode.length; i += 2) {
        uchar = (unicode[i] << 8) | unicode[i + 1]; //UNICODE为2字节编码，一次读入2个字节  
        utf8str = utf8str + String.fromCharCode(uchar); //使用String.fromCharCode强制转换  
    }
    return utf8str;
}

function ShowSpin() {
    var opts = {
        lines: 13 // The number of lines to draw
        ,
        length: 28 // The length of each line
        ,
        width: 14 // The line thickness
        ,
        radius: 42 // The radius of the inner circle
        ,
        scale: 1 // Scales overall size of the spinner
        ,
        corners: 1 // Corner roundness (0..1)
        ,
        color: '#000' // #rgb or #rrggbb or array of colors
        ,
        opacity: 0.25 // Opacity of the lines
        ,
        rotate: 0 // The rotation offset
        ,
        direction: 1 // 1: clockwise, -1: counterclockwise
        ,
        speed: 1 // Rounds per second
        ,
        trail: 60 // Afterglow percentage
        ,
        fps: 20 // Frames per second when using setTimeout() as a fallback for CSS
        ,
        zIndex: 2e9 // The z-index (defaults to 2000000000)
        ,
        className: 'spinner' // The CSS class to assign to the spinner
        ,
        top: '50%' // Top position relative to parent
        ,
        left: '50%' // Left position relative to parent
        ,
        shadow: false // Whether to render a shadow
        ,
        hwaccel: false // Whether to use hardware acceleration
        ,
        position: 'absolute' // Element positioning
    }
    if (g_spinner == null) {
        var target = document.getElementById('right_table_serverlist')
        g_spinner = new Spinner(opts).spin(target);
    }

}

function HideSpin() {
    g_spinner.spin()
    g_spinner = null
}


formatTime = function (time) {
    var unixTimestamp = new Date(time * 1000);
    return unixTimestamp.toLocaleString();
}

//需要拉取的模版列表
var Templatefilelist = ["ResourceList.html", "AddAlarmThreshold.html", "ResourceSlicesQuery.html", "AlarmList.html", "FileDirList.html", "FileList.html", "FileContent.html"]

myajax("getTemplates", JSON.stringify(Templatefilelist))

//onClickGetDirs();
function onClickGetDirs() {
    myajax("getdirs", "");
}

function GetFileList(DirName) {
    myajax("getFiles", DirName);
}

function onClickGetRoomDirs(DirName) {
    myajax("getRoomFiles", DirName);
}

function GetFileContent(FileName, levels,nPage) {
    g_levels = {"4": true, "0": true}
    g_PageNum = nPage
    g_fileName = FileName
    var tbParam = {}
    tbParam.FileName = FileName
    tbParam.Levels = levels
    tbParam.PageNum = nPage+""
    tbParam.PageSize = g_PageSize
    var strParam = JSON.stringify(tbParam)
    myajax("getFileContentWithLevel", strParam);

}

function GetFileContentByPageNum(AddOrSub,PageNum){
    var strLevels = ""
    var checkboxList = $(":checkbox")
    for (var i = 0; i < 6; i++) {
        if (g_levels[checkboxList[i].value] == true) {
            strLevels = strLevels + checkboxList[i].value + "|"
        }
    }
    var tbParam = {}
    tbParam.FileName = g_fileName
    tbParam.Levels = strLevels
    tbParam.PageSize = g_PageSize
    if(AddOrSub ==1){
        g_PageNum = g_PageNum + 1
        tbParam.PageNum = g_PageNum+""
    }
    if(AddOrSub ==-1){
        g_PageNum = g_PageNum - 1
        tbParam.PageNum = g_PageNum+""
    }
    PageNum = parseInt(PageNum)
    if(PageNum > 0){
        tbParam.PageNum = PageNum+""
        g_PageNum = PageNum
    }
    if(AddOrSub==-99){
        tbParam.PageNum = $("#PageNum").val()
        g_PageNum = parseInt(tbParam.PageNum)

    }
    var strParam = JSON.stringify(tbParam)
    myajax("getFileContentWithLevel", strParam);

}

function FilterLog(strLevel) {
    var checkboxList = $(":checkbox")
    var FileContent = {}
    FileContent.DirList = []
    g_levels = {}
    var strLevels = ""
    var bAllSelected = false
    if (strLevel == "5") {
        for (var i = 0; i < 6; i++) {
            if (checkboxList[i].value == "5" && checkboxList[i].checked) {
                bAllSelected = true;
                var tbParam = {}
                tbParam.FileName = g_fileName
                tbParam.Levels = "0|1|2|3|4|5"
                tbParam.PageNum = 1+""
                tbParam.PageSize = g_PageSize
                var strParam = JSON.stringify(tbParam)

                strLevels = "0|1|2|3|4|5"
                g_levels = {"0": true, "1": true, "2": true, "3": true, "4": true, "5": true}
                myajax("getFileContentWithLevel", strParam);
            }
        }
        return;
    }

    for (var i = 0; i < 6; i++) {
        if (checkboxList[i].value != "5" && checkboxList[i].checked == false) {
            g_levels["5"] = false;
        }
        if (checkboxList[i].checked && checkboxList[i].value != "5") {
            g_levels[checkboxList[i].value] = true
        }else{
            g_levels[checkboxList[i].value] = false
        }
    }
    for (var i = 0; i < 6; i++) {
        if (g_levels[checkboxList[i].value] == true) {
            strLevels = strLevels + checkboxList[i].value + "|"
        }
    }
    var tbParam = {}
    tbParam.FileName = g_fileName
    tbParam.Levels = strLevels
    tbParam.PageNum = 1+""
    tbParam.PageSize = g_PageSize
    var strParam = JSON.stringify(tbParam)
    myajax("getFileContentWithLevel", strParam);
    /*for (var j = 0; j < g_fileContent.DirList.length; j++) {
        if(checkboxList[5].checked==true){
            FileContent.DirList.push(g_fileContent.DirList[j])
        }else{
            for (var i = 0; i < 5; i++){
                if (checkboxList[i].checked && checkboxList[i].value == g_fileContent.DirList[j].Level) {
                    FileContent.DirList.push(g_fileContent.DirList[j])
                    break;
                }
            }
        }
    }
    $('#temp_template').html(g_all_Templates["FileContent.html"])
    var source = $('#FileContent_template').html();
    var myTemplate = Handlebars.compile(source);

    var out = myTemplate(FileContent)
    //TemplateName TemplateType LastEditTime
    $('#rightlist').html(out);
    $('#temp_template').html("");*/
}

//Plant()
//sleep(1000)
//myajax("getNodeList","")
//g_action = "getNodeList2"
//sleep(1000)  1516979233
function rnd(n, m) {
    var random = Math.floor(Math.random() * (m - n + 1) + n);
    return random;
}
