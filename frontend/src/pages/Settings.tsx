/**
 * 设置页面组件
 * 支持语言切换、年级设置等
 */

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/common/Header';
import { LanguageSwitcher } from '../components/common/LanguageSwitcher';
import { userProfileStorage } from '../services/storage';
import { gradeToAge } from '../utils/age';
import type { UserProfile } from '../types/exploration';

export default function Settings() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [grade, setGrade] = useState<string>('');

  useEffect(() => {
    loadProfile();
  }, []);


  const loadProfile = () => {
    const savedProfile = userProfileStorage.get();
    if (savedProfile) {
      setProfile(savedProfile);
      setGrade(savedProfile.grade || '');
    }
  };

  const handleGradeChange = (newGrade: string) => {
    setGrade(newGrade);
    
    // 从年级转换为年龄
    const ageFromGrade = gradeToAge(newGrade);
    const age = ageFromGrade !== undefined ? ageFromGrade : 8; // 如果转换失败，使用默认值8
    
    if (profile) {
      const updatedProfile: UserProfile = {
        ...profile,
        grade: newGrade,
        age, // 更新年龄
        lastUpdated: new Date().toISOString(),
      };
      userProfileStorage.save(updatedProfile);
      setProfile(updatedProfile);
    } else {
      // 创建新档案
      const newProfile: UserProfile = {
        age,
        grade: newGrade,
        lastUpdated: new Date().toISOString(),
      };
      userProfileStorage.save(newProfile);
      setProfile(newProfile);
    }
  };

  const grades = [
    { value: 'K1', labelKey: 'settings.gradeK1' },
    { value: 'K2', labelKey: 'settings.gradeK2' },
    { value: 'K3', labelKey: 'settings.gradeK3' },
    { value: 'G1', labelKey: 'settings.gradeG1' },
    { value: 'G2', labelKey: 'settings.gradeG2' },
    { value: 'G3', labelKey: 'settings.gradeG3' },
    { value: 'G4', labelKey: 'settings.gradeG4' },
    { value: 'G5', labelKey: 'settings.gradeG5' },
    { value: 'G6', labelKey: 'settings.gradeG6' },
    { value: 'G7', labelKey: 'settings.gradeG7' },
    { value: 'G8', labelKey: 'settings.gradeG8' },
    { value: 'G9', labelKey: 'settings.gradeG9' },
    { value: 'G10', labelKey: 'settings.gradeG10' },
    { value: 'G11', labelKey: 'settings.gradeG11' },
    { value: 'G12', labelKey: 'settings.gradeG12' },
  ];

  return (
    <div className="min-h-screen bg-cloud-white font-display">
      <Header />
      
      <main className="flex flex-col items-center px-4 py-6 w-full max-w-2xl mx-auto">
        <h1 className="text-3xl md:text-4xl font-extrabold text-slate-800 mb-8">
          {t('settings.title')}
        </h1>

        {/* 语言设置 */}
        <div className="w-full bg-white rounded-3xl border-2 border-gray-100 shadow-card p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-xl font-bold text-slate-800 mb-1">
                {t('settings.language')}
              </h2>
              <p className="text-sm text-slate-500">
                {t('settings.languageDesc')}
              </p>
            </div>
          </div>
          <LanguageSwitcher />
        </div>

        {/* 年级设置 */}
        <div className="w-full bg-white rounded-3xl border-2 border-gray-100 shadow-card p-6 mb-6">
          <div className="mb-4">
            <h2 className="text-xl font-bold text-slate-800 mb-1">
              {t('settings.grade')}
            </h2>
            <p className="text-sm text-slate-500">
              {t('settings.gradeDesc')}
            </p>
          </div>
          <div className="grid grid-cols-3 sm:grid-cols-5 gap-3">
            {grades.map((g) => (
              <button
                key={g.value}
                onClick={() => handleGradeChange(g.value)}
                className={`px-4 py-3 rounded-xl font-bold text-sm transition-all ${
                  grade === g.value
                    ? 'bg-[var(--color-primary)] text-white shadow-md'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {t(g.labelKey as any)}
              </button>
            ))}
          </div>
        </div>

        {/* 关于 */}
        <div className="w-full bg-white rounded-3xl border-2 border-gray-100 shadow-card p-6">
          <h2 className="text-xl font-bold text-slate-800 mb-4">
            {t('settings.about')}
          </h2>
          <div className="text-sm text-slate-600 space-y-2">
            <p>
              <span className="font-bold">{t('settings.version')}:</span> 1.0.0
            </p>
            <p>{t('settings.appDescription')}</p>
          </div>
        </div>

        {/* 返回按钮 */}
        <button
          onClick={() => navigate('/')}
          className="mt-8 px-8 py-3 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-600 font-bold transition-all"
        >
          {t('common.back')}
        </button>
      </main>
    </div>
  );
}
