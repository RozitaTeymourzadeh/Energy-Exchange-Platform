
let sdvalue = 0;
let fullSdValue;
let cdvalue = 0;
let fullCdValue;
const peer = '10.10.34.153';
const sdUrl='http://10.10.34.153:6686/getsdreadings';
const cdUrl='http://10.10.34.153:6686/getcdreadings';

function getSdData() {
    // value =  Math.random();
    $.ajax({
        url: sdUrl,
        crossDomain: true,
        headers: { },
        type: 'get',
        success: function (responseData) {
            // console.log(responseData[0]);
            sdvalue = responseData[0];
        },
        error: function (responseData) {
            console.log('POST failed.');
            sdvalue = 0;
        }
    });
    console.log(sdvalue);
    return sdvalue;
}

function getCdData() {
    // value =  Math.random();
    $.ajax({
        url: cdUrl,
        crossDomain: true,
        headers: { },
        type: 'get',
        success: function (responseData) {
            // console.log(responseData[0]);
            cdvalue = responseData[0];
        },
        error: function (responseData) {
            console.log('POST failed.');
            cdvalue = 0;
        }
    });
    console.log(cdvalue);
    return cdvalue;
}

var layout = {
    title:'Charge Level in : '+ peer,
    width: 1200,
    height: 600,
    xaxis: {
        showline: false,
        showgrid: false,
        showticklabels: true,
        linecolor: 'rgb(204,204,204)',
        linewidth: 2,
        autotick: true,
        ticks: 'outside',
        tickcolor: 'rgb(0,0,0)',
        tickwidth: 2,
        ticklen: 5,
        tickfont: {
            family: 'Arial',
            size: 12,
            color: 'rgb(82, 82, 82)'
        }
    },
    yaxis: {
        showgrid: false,
        zeroline: false,
        showline: false,
        showticklabels: true,
        autotick: true,
        ticks: 'outside',
        tickcolor: 'rgb(0,0,0)',
        tickwidth: 2,
        ticklen: 5,
        tickfont: {
            family: 'Arial',
            size: 12,
            color: 'rgb(82, 82, 82)'
        }
    },
    legend: {
        y: 0.5,
        // traceorder: 'reversed',
        font: {size: 18},
        yref: 'paper'
    }
};

let time = new Date();
Plotly.plot('chart',[{
        x: [time],
        y:[getSdData()],
        name: 'Supply Device',
        mode: 'lines+markers',
        line: {color: '#396AB1',width: 4},
        marker: {
            size: 12
        }
    }, {
        x: [time],
        y:[getCdData()],
        name: 'Consume Device',
        mode: 'lines+markers',
        line: {color: '#948B3D',width: 4},
        marker: {
            size: 12
        }
    }],
    layout,
    {responsive: true});


var cnt = 0;
setInterval(function (){
    let time = new Date();
    let update = {
        x:  [[time],[time]],
        y: [[getSdData()],[getCdData()]]
    };

    Plotly.extendTraces('chart', update, [0,1])

    cnt++;

    if(cnt > 500) {
        Plotly.relayout('chart',{
            xaxis: {
                range: [cnt-500,cnt]
            }
        });
    }
},8000);