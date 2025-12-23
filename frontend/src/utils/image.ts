/**
 * 图片处理工具函数
 */

/**
 * 将File对象转换为base64字符串
 */
export async function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result as string;
      resolve(result);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

/**
 * 压缩图片
 * @param file 原始图片文件
 * @param maxWidth 最大宽度
 * @param maxHeight 最大高度
 * @param quality 压缩质量 0-1
 * @returns 压缩后的 Blob（JPEG 格式）
 */
export async function compressImage(
  file: File,
  maxWidth: number = 1920,
  maxHeight: number = 1920,
  quality: number = 0.8
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => {
      const canvas = document.createElement('canvas');
      let width = img.width;
      let height = img.height;

      // 计算缩放比例
      if (width > maxWidth || height > maxHeight) {
        const ratio = Math.min(maxWidth / width, maxHeight / height);
        width = width * ratio;
        height = height * ratio;
      }

      canvas.width = width;
      canvas.height = height;

      const ctx = canvas.getContext('2d');
      if (!ctx) {
        reject(new Error('无法创建canvas上下文'));
        return;
      }

      ctx.drawImage(img, 0, 0, width, height);

      // 统一转换为 JPEG 格式
      canvas.toBlob(
        (blob) => {
          if (!blob) {
            reject(new Error('图片压缩失败'));
            return;
          }
          resolve(blob);
        },
        'image/jpeg', // 统一使用 JPEG 格式
        quality
      );
    };
    img.onerror = reject;
    img.src = URL.createObjectURL(file);
  });
}

/**
 * 清理base64字符串，移除所有空白字符
 * 这可以解决传输过程中可能引入的空格、换行符等问题
 */
export function cleanBase64String(s: string): string {
  // 移除所有空白字符（空格、换行符、制表符等）
  return s.replace(/\s/g, '');
}

/**
 * 从base64字符串中提取数据部分（去除data:image/...;base64,前缀）
 * 并清理可能的空白字符（防御性编程）
 */
export function extractBase64Data(base64: string): string {
  const commaIndex = base64.indexOf(',');
  const data = commaIndex > 0 ? base64.substring(commaIndex + 1) : base64;
  // 清理可能的空白字符，确保base64字符串干净
  return cleanBase64String(data);
}

