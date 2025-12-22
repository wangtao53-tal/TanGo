/**
 * 年龄相关工具函数
 * 处理年级到年龄的转换
 */

import type { UserProfile } from '../types/exploration';
import { userProfileStorage } from '../services/storage';

/**
 * 年级到年龄的映射
 * K1-K3: 幼儿园（3-6岁）
 * G1-G6: 小学（7-12岁）
 * G7-G9: 初中（13-15岁）
 * G10-G12: 高中（16-18岁）
 */
const GRADE_TO_AGE_MAP: Record<string, number> = {
  'K1': 4,  // 幼儿园小班 -> 4岁
  'K2': 5,  // 幼儿园中班 -> 5岁
  'K3': 6,  // 幼儿园大班 -> 6岁
  'G1': 7,  // 一年级 -> 7岁
  'G2': 8,  // 二年级 -> 8岁
  'G3': 9,  // 三年级 -> 9岁
  'G4': 10, // 四年级 -> 10岁
  'G5': 11, // 五年级 -> 11岁
  'G6': 12, // 六年级 -> 12岁
  'G7': 13, // 七年级 -> 13岁
  'G8': 14, // 八年级 -> 14岁
  'G9': 15, // 九年级 -> 15岁
  'G10': 16, // 十年级 -> 16岁
  'G11': 17, // 十一年级 -> 17岁
  'G12': 18, // 十二年级 -> 18岁
};

/**
 * 将年级转换为年龄
 * @param grade 年级（K1-K12格式）
 * @returns 年龄（3-18岁），如果年级无效则返回undefined
 */
export function gradeToAge(grade?: string): number | undefined {
  if (!grade) {
    return undefined;
  }
  return GRADE_TO_AGE_MAP[grade];
}

/**
 * 从用户档案获取年龄
 * 优先级：年级转换 > 存储的年龄 > 默认值8
 * @param profile 用户档案（可选）
 * @param defaultAge 默认年龄，默认为8
 * @returns 年龄（3-18岁）
 */
export function getUserAge(profile?: UserProfile | null, defaultAge: number = 8): number {
  // 1. 优先从年级转换
  if (profile?.grade) {
    const ageFromGrade = gradeToAge(profile.grade);
    if (ageFromGrade !== undefined) {
      return ageFromGrade;
    }
  }

  // 2. 使用存储的年龄
  if (profile?.age && profile.age >= 3 && profile.age <= 18) {
    return profile.age;
  }

  // 3. 使用默认值
  return defaultAge;
}

/**
 * 从存储中获取用户年龄
 * 从 userProfileStorage 读取并转换为年龄
 * @param defaultAge 默认年龄，默认为8
 * @returns 年龄（3-18岁）
 */
export function getUserAgeFromStorage(defaultAge: number = 8): number {
  const profile = userProfileStorage.get();
  return getUserAge(profile, defaultAge);
}

