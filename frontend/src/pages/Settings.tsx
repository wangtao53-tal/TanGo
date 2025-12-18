/**
 * 设置页面组件
 * 支持语言切换、年级设置等
 */

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/common/Header';
import { LanguageSwitcher } from '../components/common/LanguageSwitcher';
import { userSettingsStorage, userProfileStorage } from '../services/storage';
import type { UserSettings } from '../types/settings';
import type { UserProfile } from '../types/exploration';

export default function Settings() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [settings, setSettings] = useState<UserSettings | null>(null);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [grade, setGrade] = useState<string>('');

  useEffect(() => {
    loadSettings();
    loadProfile();
  }, []);

  const loadSettings = () => {
    const savedSettings = userSettingsStorage.get();
    if (savedSettings) {
      setSettings(savedSettings);
    } else {
      const defaultSettings = userSettingsStorage.getDefault();
      userSettingsStorage.save(defaultSettings);
      setSettings(defaultSettings);
    }
  };

  const loadProfile = () => {
    const savedProfile = userProfileStorage.get();
    if (savedProfile) {
      setProfile(savedProfile);
      setGrade(savedProfile.grade || '');
    }
  };

  const handleGradeChange = (newGrade: string) => {
    setGrade(newGrade);
    if (profile) {
      const updatedProfile: UserProfile = {
        ...profile,
        grade: newGrade,
        lastUpdated: new Date().toISOString(),
      };
      userProfileStorage.save(updatedProfile);
      setProfile(updatedProfile);
    } else {
      // 创建新档案
      const newProfile: UserProfile = {
        age: 8, // 默认年龄
        grade: newGrade,
        lastUpdated: new Date().toISOString(),
      };
      userProfileStorage.save(newProfile);
      setProfile(newProfile);
    }
  };

  const grades = [
    { value: 'K1', label: 'Kindergarten 1' },
    { value: 'K2', label: 'Kindergarten 2' },
    { value: 'K3', label: 'Kindergarten 3' },
    { value: 'G1', label: 'Grade 1' },
    { value: 'G2', label: 'Grade 2' },
    { value: 'G3', label: 'Grade 3' },
    { value: 'G4', label: 'Grade 4' },
    { value: 'G5', label: 'Grade 5' },
    { value: 'G6', label: 'Grade 6' },
    { value: 'G7', label: 'Grade 7' },
    { value: 'G8', label: 'Grade 8' },
    { value: 'G9', label: 'Grade 9' },
    { value: 'G10', label: 'Grade 10' },
    { value: 'G11', label: 'Grade 11' },
    { value: 'G12', label: 'Grade 12' },
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
                {g.label}
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
            <p>TanGo - 探索世界的知识卡片应用</p>
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
