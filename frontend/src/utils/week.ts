/**
 * 周相关工具函数
 */

/**
 * 获取当前周的开始时间（周一 00:00:00）和结束时间（周日 23:59:59）
 * @returns { start: Date, end: Date } 当前周的开始和结束时间
 */
export function getCurrentWeekRange(): { start: Date; end: Date } {
  const now = new Date();
  const day = now.getDay(); // 0 = 周日, 1 = 周一, ..., 6 = 周六
  
  // 计算到本周一的偏移天数
  // 如果今天是周日(0)，则偏移-6天；如果是周一(1)，则偏移0天；以此类推
  const offsetToMonday = day === 0 ? -6 : 1 - day;
  
  // 本周一的开始时间（00:00:00）
  const monday = new Date(now);
  monday.setDate(now.getDate() + offsetToMonday);
  monday.setHours(0, 0, 0, 0);
  
  // 本周日的结束时间（23:59:59.999）
  const sunday = new Date(monday);
  sunday.setDate(monday.getDate() + 6);
  sunday.setHours(23, 59, 59, 999);
  
  return {
    start: monday,
    end: sunday,
  };
}

/**
 * 判断一个时间戳是否在当前周内
 * @param timestamp ISO 8601 格式的时间戳字符串
 * @returns 是否在当前周内
 */
export function isInCurrentWeek(timestamp: string): boolean {
  const { start, end } = getCurrentWeekRange();
  const date = new Date(timestamp);
  return date >= start && date <= end;
}

/**
 * 获取当前周的日期范围字符串（用于显示）
 * @returns 格式化的日期范围字符串，如 "2025/12/22 - 2025/12/28"
 */
export function getCurrentWeekRangeString(): string {
  const { start, end } = getCurrentWeekRange();
  const formatDate = (date: Date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}/${month}/${day}`;
  };
  return `${formatDate(start)} - ${formatDate(end)}`;
}

