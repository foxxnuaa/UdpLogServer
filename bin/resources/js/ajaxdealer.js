dealWithAjaxRetData = function(data, status) {
    //getTemplates
    g_oldRightList = $('#rightlist').html()
    if ("getTemplates" == g_action) {
        g_all_Templates = JSON.parse(data)
        g_all_Templates = JSON.parse(g_all_Templates.data)
    }

    if("recourcelist" == g_action){
        //var g_resourcelist= {}
        //var g_resourceid = 1
        g_resourceid = 1
        g_resourcelist = JSON.parse(JSON.parse(data)["data"])
        var ResourceList = {}
        ResourceList.ResourceList = []
        var TypeList = []
        var Idlist = []
        if (g_resourcelist != null){
            for (var i = g_resourcelist.length - 1; i >= 0; i--){
                if(g_resourcelist[i].FiveMinuteThreshold > 0 || g_resourcelist[i].OneHourThreshold > 0 ||g_resourcelist[i].OneDayThreshold > 0){
                    ResourceList.ResourceList.push(g_resourcelist[i])
                }
                g_resourcelist[i].resourceid =g_resourceid
                g_resourceid = g_resourceid + 1
                TypeList.push(g_resourcelist[i].Type)
                Idlist.push(g_resourcelist[i].Id)
            }
        }

        $('#temp_template').html(g_all_Templates["ResourceList.html"])

        var source = $('#ResourceList-template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(ResourceList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
    }

    if("onclickaddalarm" == g_action){
        //var g_resourcelist= {}
        //var g_resourceid = 1
        g_resourceid = 1
        g_resourcelist = JSON.parse(JSON.parse(data)["data"])
        var ResourceList = {}
        ResourceList.ResourceList = []
        var TypeList = []
        var Idlist = []
        if (g_resourcelist != null){
            for (var i = g_resourcelist.length - 1; i >= 0; i--){
                if(g_resourcelist[i].FiveMinuteThreshold == 0 || g_resourcelist[i].OneHourThreshold == 0 ||g_resourcelist[i].OneDayThreshold == 0){
                    ResourceList.ResourceList.push(g_resourcelist[i])
                }
                g_resourcelist[i].resourceid =g_resourceid
                g_resourceid = g_resourceid + 1
            }
        }

        $('#temp_template').html(g_all_Templates["ResourceList.html"])

        var source = $('#ResourceList-template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(ResourceList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
    }
    if("OnClickShowResourceSlices" ==g_action){
        var resourceSlices = JSON.parse(JSON.parse(data)["data"])
        var ResourceList = {}
        ResourceList.ResourceList = resourceSlices
        $('#temp_template').html(g_all_Templates["ResourceSlicesQuery.html"])

        var source = $('#ResourceSlice-template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(ResourceList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        $('#temp_template').html("");
        $('.multiselect').multiselect();
        $('#example28').multiselect({
            includeSelectAllOption: true,
            enableFiltering: true,
            maxHeight: 150
        });
        Datetime()
    }
    if("getResourceSlice" == g_action){
        var resourceSlices = JSON.parse(JSON.parse(data)["data"])
        var ResourceList = [[],[]]
        var strDay = document.getElementById("nowdate").value
        strDay =  strDay.replace(/-/g,"")
        for(var key in resourceSlices){
            if (key==strDay){
                ResourceList[0] = JSON.parse(resourceSlices[key])
            }else{
                ResourceList[1] = resourceSlices[key]
            }
        }
        Plant(ResourceList)
        //Plant(AllSlices)
    }
    if("getdirs"==g_action){
        var DirList1 = JSON.parse(JSON.parse(data)["data"])
        var DirList = {}
        DirList.DirList = DirList1
        $('#temp_template').html(g_all_Templates["FileDirList.html"])

        var source = $('#DirList_template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(DirList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        $('#temp_template').html("");
    }
    if("getFiles"==g_action){
        var DirList1 = JSON.parse(JSON.parse(data)["data"])
        var DirList = {}
        DirList.DirList = DirList1
        $('#temp_template').html(g_all_Templates["FileList.html"])

        var source = $('#FileList_template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(DirList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        $('#temp_template').html("");
    }

    if("getRoomFiles" == g_action){
        var DirList1 = JSON.parse(JSON.parse(data)["data"])
        var DirList = {}
        DirList.DirList = DirList1
        var compare = function (obj1, obj2) {
            var val1 = parseInt(obj1.FileName);
            var val2 = parseInt(obj2.FileName);
            if (val1 > val2) {
                return -1;
            } else if (val1 < val2) {
                return 1;
            } else {
                return 0;
            }
        }
        DirList.DirList.sort(compare)
        $('#temp_template').html(g_all_Templates["FileList.html"])

        var source = $('#FileList_template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(DirList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        $('#temp_template').html("");
    }

    if("getFileContent"==g_action||"getFileContentWithLevel" == g_action){
        var DirList1 = JSON.parse(JSON.parse(data)["data"])
        g_fileContent.DirList = JSON.parse(DirList1.FileContent)
        g_fileContent.TotalPageNum = DirList1.TotalPageNum
        g_fileContent.CurrentPageNum = g_PageNum    
        var CalcPageNumsResult = CalcPageNums(g_PageNum,DirList1.TotalPageNum)
        g_fileContent.PageItems = CalcPageNumsResult.PageItems

        $('#temp_template').html(g_all_Templates["FileContent.html"])

        var source = $('#FileContent_template').html();
        var myTemplate = Handlebars.compile(source);

        var out = myTemplate(g_fileContent)
        out = out.replace(/\*\*&amp/g,"<br/>")
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);

        var checkboxList = $(":checkbox")
            for (var j = 0; j < checkboxList.length; j++) {
                if(g_levels[checkboxList[j].value] == true){
                    checkboxList[j].checked = true
                }else{
                    checkboxList[j].checked = false
                }
            }
        $('#temp_template').html("");
        $("#PageItem_" + g_PageNum).addClass("active")
        if(parseInt(g_fileContent.CurrentPageNum) >= parseInt(g_fileContent.TotalPageNum))
        {
            $("#PageItem_Next").addClass("disabled")
            $("#PageItem_Next").children().removeAttr("onclick");
        }
        if(parseInt(g_fileContent.CurrentPageNum) == 1)
        {
            $("#PageItem_Pre").addClass("disabled")
            $("#PageItem_Pre").children().removeAttr("onclick");
        }
    }
    //alert("Data: " + data + "\nStatus: " + status);
    history.pushState({}, "兵锋日志系统","#"+g_action);

    $(window).on("popstate",function(){
        $('#rightlist').html(g_oldRightList);
    });
}

CalcPageNums = function(CurrentPageNo, TotalPageNum) {
    var ret = {}
    ret.PageItems = []
    if(TotalPageNum >= 10)
    {
        if(CurrentPageNo <= 10)
        {
           ret.BeginPageNo = 1
           ret.EndPageNo = 10
        }
        else
        {
            ret.BeginPageNo = CurrentPageNo - 4
            ret.EndPageNo = (CurrentPageNo+5)>TotalPageNum?TotalPageNum:CurrentPageNo+5
        }       
    } 
    else
    {
        ret.BeginPageNo = 1
        ret.EndPageNo = TotalPageNum
    }
    for (var i = ret.BeginPageNo; i <= ret.EndPageNo; i++)
    {
        ret.PageItems.push(i)
    }
    return ret
}