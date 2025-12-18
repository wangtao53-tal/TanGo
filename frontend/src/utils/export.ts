/**
 * 导出工具函数
 * 使用 html2canvas 将卡片导出为图片
 */

import html2canvas from 'html2canvas';

/**
 * 导出卡片为图片
 * @param elementId 卡片元素的ID
 * @param filename 文件名（不含扩展名）
 */
export async function exportCardAsImage(
  elementId: string,
  filename: string = 'card'
): Promise<void> {
  const element = document.getElementById(elementId);
  if (!element) {
    throw new Error(`找不到ID为 ${elementId} 的元素`);
  }

  try {
    // 使用 html2canvas 捕获元素
    const canvas = await html2canvas(element, {
      backgroundColor: '#ffffff',
      scale: 2, // 提高图片清晰度
      useCORS: true,
      logging: false,
    });

    // 转换为 blob
    canvas.toBlob((blob) => {
      if (!blob) {
        throw new Error('生成图片失败');
      }

      // 创建下载链接
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${filename}.png`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    }, 'image/png');
  } catch (error) {
    console.error('导出卡片失败:', error);
    throw error;
  }
}

/**
 * 导出多个卡片为图片（批量导出）
 * @param elementIds 卡片元素ID数组
 * @param baseFilename 基础文件名
 */
export async function exportCardsAsImages(
  elementIds: string[],
  baseFilename: string = 'cards'
): Promise<void> {
  for (let i = 0; i < elementIds.length; i++) {
    const elementId = elementIds[i];
    await exportCardAsImage(elementId, `${baseFilename}-${i + 1}`);
    // 添加延迟，避免浏览器阻止多个下载
    if (i < elementIds.length - 1) {
      await new Promise((resolve) => setTimeout(resolve, 500));
    }
  }
}
