// Hopefully these utility functions will not be needed in the future

export function parsePercentData(data) {
  if (!data) {
    return [];
  }
  const newData = [...data];
  newData.forEach((element) => {
    if (element.name || element.y) {
      return element;
    }
    element.name = element.key;
    delete element.key;
    element.y = Math.round(element.value[0].value * 100) / 100;
    delete element.value;
  });
  newData.sort(function (a, b) {
    return b.y - a.y;
  });
  const output = newData.filter((item) => item.name !== "ALL");

  return output;
}

export function parseDistrictData(data, index) {
  const newData = data.filter((item) => item.key !== "e");
  const differentDataNum = newData[0].value.length;
  const xAxis = newData.map((item) => {
    if (item.key === "e") {
      return;
    }
    return item.key;
  });
  const series = [];
  const district = {};

  console.log(newData);

  for (let i = 0; i < differentDataNum; i++) {
    let dataObj = { name: newData[0].value[i].key, data: [] };
    for (let j = 0; j < newData.length; j++) {
      dataObj.data.push(+newData[j].value[i].value);
    }
    series.push(dataObj);
  }

  for (let i = 0; i < newData.length; i++) {
    for (let j = 0; j < differentDataNum; j++) {
      if (i === 0 && j === 0) {
        district.highest = 0;
      }
      if (j === 0) {
        district[newData[i].key] = 0;
      }
      district[newData[i].key] += +newData[i].value[j].value;
      if (
        j === differentDataNum - 1 &&
        district[newData[i].key] > district.highest
      ) {
        district.highest = district[newData[i].key];
      }
    }
  }
  return [xAxis, series, district];
}

export function sumAllData(data) {
  let sum = 0;
  for (let i = 0; i < data.length; i++) {
    sum += data[i].y;
  }
  return Math.round(sum * 100) / 100;
}

export function parseAreaData(data, color) {
  if (!data) {
    return [];
  }
  let newData = [...data];
  newData.forEach((element) => {
    if (element.name || element.y) {
      return element;
    }
    element.name = element.key;
    delete element.key;
    element.y = Math.round(element.value[0].value * 100) / 100;
    element.value = Math.round(element.value[0].value * 100) / 100;
  });
  if (color) {
    newData.forEach((element) => {
      element.color = color[element.name];
    });
  }
  newData = newData.filter(
    (item) => !["國有", "臺北市有", "台北市總面積"].includes(item.name)
  );
  newData.sort(function (a, b) {
    return b.y - a.y;
  });
  const output = newData.filter((item) => item.name !== "ALL");

  return output;
}

export function parseTimeData(data, time) {
  if (!data) {
    return [];
  }
  let newData = [...data];
  newData.forEach((element) => {
    if (element.name || element.data) {
      return element;
    }
    element.name = element.key;
    delete element.key;
    const extractedData = [];
    for (let i = 0; i < element.value.length; i++) {
      extractedData.push(element.value[i].value);
    }
    element.data = extractedData;
    delete element.value;
  });

  newData = newData.filter((item) => !["破獲率[%]"].includes(item.name));

  console.log(newData);

  return [newData, time];
}

export function parseGuageData(data) {
  return data;
}
