// 這是小龜的魔法日期轉換程式，目前我還沒有看懂但他似乎包含了所有狀況了，所以就這樣吧

import dayjs from "dayjs";

const DateTimeFormat = {
  default: "YYYY-MM-DDTHH:mm:ssZ",
  initTime: "YYYY-MM-DDT00:00:00Z",
  initHour: "YYYY-MM-DDTHH:00:00Z",
  lastTime: "YYYY-MM-DDT23:59:59Z",
  overnightTime: "YYYY-MM-DDT12:00:00Z",
};

export const parseTimeFun = (dataObj) => {
  /**
   * request parameter
   * last: target last datetime
   * time_step: -1, 0, 1+
   * time_unit: minute | hour | day  | week | month
   * stay_mins: 停留時間(minute)
   */
  let { last, time_step, time_unit } = { ...dataObj };
  /**
   * Response data
   * start: start time
   * end: end time
   */
  let newObj = {};
  if (typeof time_step == "undefined" && typeof time_unit == "undefined")
    return {};
  if (last) {
    newObj.end = dayjs(last).format(DateTimeFormat.default);
  } else {
    newObj.end = dayjs();
  }
  if (time_step == 0) {
    return parseCustomTime(newObj, time_unit);
  } else if (time_step == 1) {
    return parsePartTime(newObj, time_unit);
  }
  newObj.start = dayjs(newObj.end).subtract(time_step, time_unit);
  if (time_unit === "hour") {
    newObj.start = dayjs(newObj.start).format("YYYY-MM-DDTHH:00:00Z");
  } else if (time_unit === "month") {
    newObj.start = dayjs(newObj.start).date(1).format(DateTimeFormat.initTime);
  } else {
    newObj.start = dayjs(newObj.start).format(DateTimeFormat.default);
  }
  return {
    start: newObj.start,
    end: dayjs(newObj.end).format(DateTimeFormat.default),
  };
};

const parseCustomTime = (customObj, time_unit) => {
  if (time_unit === "week") {
    //上周:上週日-上週六共七日
    customObj.end = dayjs(customObj.end).day(-1).format(DateTimeFormat.default);
    customObj.start = dayjs(customObj.end)
      .subtract(6, "day")
      .format(DateTimeFormat.initTime);
  } else if (time_unit === "day") {
    //今日12:00-隔日12:00
    const isOvernight =
      dayjs(customObj.end).hour() >= 0 && dayjs(customObj.end).hour() < 13;
    if (isOvernight) {
      // 現在時間為跨夜後00:00-11:59 -> 請求時間：前天12:00
      customObj.start = dayjs(customObj.end)
        .day(-1)
        .format(DateTimeFormat.overnightTime);
    } else {
      // 現在時間為跨夜前12:00-23:59 -> 請求時間：今天12:00
      customObj.start = dayjs(customObj.end).format(
        DateTimeFormat.overnightTime
      );
    }
  }
  return customObj;
};

const parsePartTime = (customObj, time_unit) => {
  if (time_unit === "week") {
    //當周:上週日開始計算
    const isSunday = dayjs(customObj.end).day() === 0;
    if (isSunday) {
      customObj.start = dayjs(customObj.end)
        .subtract(7, "day")
        .format(DateTimeFormat.initTime);
    } else {
      customObj.start = dayjs(customObj.end)
        .day(0)
        .format(DateTimeFormat.initTime);
    }
  } else if (time_unit === "hour") {
    // For MRT
    customObj = parseMinFormat(customObj.end);
  } else if (time_unit === "day") {
    //今日00:00-隔日23:59
    customObj.start = dayjs(customObj.end).format(DateTimeFormat.initTime);
  } else if (time_unit === "month") {
    customObj.start = dayjs(customObj.end)
      .date(1)
      .format(DateTimeFormat.initTime);
  } else {
    customObj.start = dayjs(customObj.end)
      .subtract(1, time_unit)
      .format(DateTimeFormat.initTime);
  }
  customObj.end = dayjs(customObj.end).format(DateTimeFormat.default);
  return customObj;
};

const parseMinFormat = (end) => {
  // Data is stored once in 30 mins，only "XX:29, XX: 59"
  // Insert into database delayed 5mins
  const minutes = dayjs().minute();
  const newObj = {};
  if (minutes < 5) {
    // current: 02:00-02:04
    // get: 01:29
    // request: 01:00-01: 30
    newObj.start = dayjs(end)
      .subtract(1, "hour")
      .format("YYYY-MM-DDTHH:00:00Z");
    newObj.end = dayjs(end).subtract(1, "hour").format("YYYY-MM-DDTHH:30:00Z");
  } else if (minutes >= 5 && minutes < 35) {
    // current: 02:05-02:34
    // get: 01:59
    // request:01:31-02: 00
    newObj.start = dayjs(end)
      .subtract(1, "hour")
      .format("YYYY-MM-DDTHH:31:00Z");
    newObj.end = dayjs(end).format("YYYY-MM-DDTHH:00:00Z");
  } else {
    // current: 02:35-02:59
    // get: 02:29
    // request: 02:00-02: 30
    newObj.start = dayjs(end).format("YYYY-MM-DDTHH:00:00Z");
    newObj.end = dayjs(end).format("YYYY-MM-DDTHH:30:00Z");
  }
  return newObj;
};
