// js part of the /pool/:num
// this scripts gets chart data from html hidden elements
// and display it on chart via chartJS library.

var ctx = document.getElementById("resultsChart");
var chartData = document.getElementById('chart-data');

var votes = [];
var labels = [];

// getting data out of hidden html elements
var data = chartData.children;

for (var i = 0; i < data.length; i++) {
    if (i % 2 == 0) { // get titles of our data
        labels.push(data[i].defaultValue);
    } else { // push votes of our data
        votes.push(data[i].defaultValue);
    }
}

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
                    beginAtZero: true,
                    //removing decimal points from table
                    userCallback: function (label, index, labels) {
                        if (Math.floor(label) === label) {
                            return label;
                        }
                    }
                }
            }]
        }
    }
});