// js part of the /poll/:num

// this scripts gets number of votes for each vote  option
//from html hidden elements and display it on chart
//via chartJS library.

var ctx = document.getElementById("resultsChart");
var chartData = document.getElementById('chart-data');

var votes = [];
var labels = [];

// parsing data out of hidden html elements that are filled
// via go templates
var data = chartData.children;

for (var i = 0; i < data.length; i++) {
    if (i % 2 == 0) { // get titles of our data
        labels.push(data[i].defaultValue);
    } else { // push votes of our data
        votes.push(data[i].defaultValue);
    }
}


// getMaxValue(["1", "2", "4"]) returns a number that is used as the maximum number
// to defines maximum height of the chart. Setting max number on chart ensures
//our yAxis scales nicely to the nearest suitable number: See below
//
// Algorithm:
// 53 turned 60 => nearest 2 digit number
// 120 turned into 200 => nearest 3 digit number
// 1200 turned into 2000 => nearest 4 digit number
//
function getMaxValue(valuesArr) {
    // turn array of strings into array of integers
    arr = valuesArr.map(function (item) {
        return parseInt(item);
    })
    // get maximum number in array
    max = Math.max.apply(null, arr)

    // algorithm to turn numbers into nearest suitable number
    // see description above this function
    var padding = ("" + max).length;
    var scale = 10 ** (padding - 1);
    var maxValue = Math.floor(max / scale) * scale + scale;
    return maxValue;
}


var maxValue = getMaxValue(votes)

var myChart = new Chart(ctx, {
    type: 'bar',
    data: {
        labels: labels,
        datasets: [{
            //label: '# of Votes',
            data: votes,
            backgroundColor: [
                'rgba(255, 99, 132, 0.2)',
                'rgba(54, 162, 235, 0.2)',
                'rgba(255, 206, 86, 0.2)',
                'rgba(75, 192, 192, 0.2)',
                'rgba(153, 102, 255, 0.2)',
                'rgba(255, 159, 64, 0.2)'
            ],
            borderColor: [
                'rgba(255,99,132,1)',
                'rgba(54, 162, 235, 1)',
                'rgba(255, 206, 86, 1)',
                'rgba(75, 192, 192, 1)',
                'rgba(153, 102, 255, 1)',
                'rgba(255, 159, 64, 1)'
            ],
            borderWidth: 1
        }]
    },
    options: {
        responsive: true,
        maintainAspectRatio: true,
        legend: {
            display: false,
        },
        scales: {
            yAxes: [{
                scaleLabel: {
                    display: true,
                    labelString: "Number of Votes"
                },
                ticks: {
                    padding: 5,
                    max: maxValue,
                    beginAtZero: true,
                    //removing decimal points from table
                    userCallback: function (label, index, labels) {
                        if (Math.floor(label) === label) {
                            return label;
                        }
                    },
                }
            }]
        },
    }
});

// This plugin draws numbers above chart bars, taken from stackoverflow
Chart.plugins.register({
    afterDatasetsDraw: function (chart, easing) {
        // To only draw at the end of animation, check for easing === 1
        var ctx = chart.ctx;
        chart.data.datasets.forEach(function (dataset, i) {
            var meta = chart.getDatasetMeta(i);
            if (!meta.hidden) {
                meta.data.forEach(function (element, index) {
                    // Draw the text in black, with the specified font
                    ctx.fillStyle = 'rgb(0, 0, 0)';
                    var fontSize = 16;
                    var fontStyle = 'normal';
                    var fontFamily = 'Fira Sans';
                    ctx.font = Chart.helpers.fontString(fontSize, fontStyle, fontFamily);
                    // Just naively convert to string for now
                    var dataString = dataset.data[index].toString();
                    // Make sure alignment settings are correct
                    ctx.textAlign = 'center';
                    ctx.textBaseline = 'middle';
                    var padding = 0;
                    var position = element.tooltipPosition();
                    ctx.fillText(dataString, position.x, position.y - (fontSize / 2) - padding);
                });
            }
        });
    }
});