import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { useEffect } from 'react';
import Home from './pages/Home';
import Capture from './pages/Capture';
import Result from './pages/Result';
import Collection from './pages/Collection';
import Share from './pages/Share';
import LearningReport from './pages/LearningReport';
import Settings from './pages/Settings';
import Badge from './pages/Badge';
import { QuickCaptureButton } from './components/common/QuickCaptureButton';
import { userSettingsStorage } from './services/storage';
import './i18n'; // 初始化i18n

function App() {
  // 应用启动时加载语言设置
  useEffect(() => {
    const settings = userSettingsStorage.get();
    if (settings) {
      // i18n会在初始化时自动从localStorage读取语言设置
      // 这里确保设置已保存
      if (!settings.lastUpdated) {
        userSettingsStorage.save({
          ...settings,
          lastUpdated: new Date().toISOString(),
        });
      }
    } else {
      // 如果没有设置，使用默认设置
      const defaultSettings = userSettingsStorage.getDefault();
      userSettingsStorage.save(defaultSettings);
    }
  }, []);

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-cloud-white">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/capture" element={<Capture />} />
          <Route path="/result" element={<Result />} />
          <Route path="/collection" element={<Collection />} />
          <Route path="/share/:shareId" element={<Share />} />
          <Route path="/report" element={<LearningReport />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="/badge" element={<Badge />} />
        </Routes>
        {/* 全局快速拍照按钮 */}
        <QuickCaptureButton />
      </div>
    </BrowserRouter>
  );
}

export default App;
