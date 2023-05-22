export const Options = {
  title: {
    text: null,
    style: {
      display: "none",
    },
  },
  credits: {
    enabled: false,
  },
};

export const ChartConfig = {
  marginTop: -10,
  spacingTop: -10,
};

export const LineXAxis = {
  // lineWidth: null,
  lineColor: "rgba(200,200,200,0.2)",
  gridLineWidth: 0,
  minPadding: 0,
  maxPadding: 0,
};

export const LineYAxis = {
  title: {
    text: null,
  },
  lineWidth: 1,
  lineColor: "rgba(200,200,200,0.2)",
  // gridLineWidth: 0,
  gridLineColor: "rgba(100,100,100,0.1)",
  minorGridLineWidth: 0,
  min: 0,
  // tickWidth: 0,
  // minRange: 0.001,
  labels: {
    formatter: function () {
      // return this.value / 1000 + 'k';
      return this.value;
    },
  },
};
