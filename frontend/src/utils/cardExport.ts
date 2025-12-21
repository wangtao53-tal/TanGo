/**
 * 卡片导出工具
 * 封装 html2canvas，用于将卡片导出为图片
 */

import html2canvas from 'html2canvas';

export interface CardExportOptions {
  scale?: number; // 图片清晰度（默认3）
  backgroundColor?: string; // 背景色（默认白色）
  filename?: string; // 文件名（不含扩展名）
}

/**
 * 导出卡片为图片
 * @param elementId 卡片元素的ID
 * @param options 导出选项
 */
export async function exportCardAsImage(
  elementId: string,
  options: CardExportOptions = {}
): Promise<void> {
  const {
    scale = 3,
    backgroundColor = '#ffffff',
    filename = `card-${Date.now()}`,
  } = options;

  const element = document.getElementById(elementId);
  if (!element) {
    throw new Error(`找不到ID为 ${elementId} 的元素`);
  }

  try {
    // 使用 html2canvas 捕获元素，优化配置确保清晰度
    const canvas = await html2canvas(element, {
      backgroundColor,
      scale, // 提高图片清晰度
      useCORS: true,
      allowTaint: false,
      logging: false,
      width: element.scrollWidth,
      height: element.scrollHeight,
      windowWidth: element.scrollWidth,
      windowHeight: element.scrollHeight,
      // 确保所有样式都被渲染
      onclone: (clonedDoc, element) => {
        // 第一步：移除或禁用所有包含 oklch 的样式表规则
        // 这是最关键的步骤，必须在 html2canvas 解析样式之前完成
        try {
          const styleSheets = clonedDoc.styleSheets;
          const styleSheetsArray = Array.from(styleSheets);
          
          styleSheetsArray.forEach((sheet) => {
            try {
              const rules = sheet.cssRules || sheet.rules;
              if (!rules) return;
              
              const rulesArray = Array.from(rules);
              // 从后往前删除，避免索引问题
              for (let i = rulesArray.length - 1; i >= 0; i--) {
                const rule = rulesArray[i];
                try {
                  const cssText = rule.cssText || '';
                  if (cssText.includes('oklch')) {
                    sheet.deleteRule(i);
                  } else if (rule instanceof CSSMediaRule || rule instanceof CSSKeyframesRule) {
                    // 处理嵌套规则（@media、@keyframes 等）
                    const nestedRules = rule.cssRules;
                    if (nestedRules) {
                      const nestedArray = Array.from(nestedRules);
                      for (let j = nestedArray.length - 1; j >= 0; j--) {
                        const nestedRule = nestedArray[j];
                        if (nestedRule.cssText && nestedRule.cssText.includes('oklch')) {
                          rule.deleteRule(j);
                        }
                      }
                    }
                  }
                } catch (e) {
                  // 忽略无法删除的规则
                }
              }
            } catch (e) {
              // 跨域样式表可能无法访问，忽略
              console.warn('无法访问样式表:', e);
            }
          });
        } catch (e) {
          console.warn('处理样式表时出错:', e);
        }

        const clonedElement = clonedDoc.getElementById(elementId);
        if (!clonedElement || !clonedDoc.defaultView) return;

        // 第二步：遍历所有元素，清理 style 属性中的 oklch，并使用 getComputedStyle 设置 RGB 值
        const allElements = clonedElement.querySelectorAll('*');
        const elementsArray = [clonedElement, ...Array.from(allElements)];
        
        elementsArray.forEach((el) => {
          const htmlEl = el as HTMLElement;
          if (!htmlEl) return;

          try {
            // 清理 style 属性中的 oklch
            const inlineStyle = htmlEl.getAttribute('style');
            if (inlineStyle && inlineStyle.includes('oklch')) {
              // 移除包含 oklch 的样式属性
              const styleParts = inlineStyle.split(';').filter(part => {
                return part.trim() && !part.includes('oklch');
              });
              htmlEl.setAttribute('style', styleParts.join(';'));
            }

            // 使用 getComputedStyle 获取计算后的 RGB 值并设置为内联样式
            const computedStyle = clonedDoc.defaultView.getComputedStyle(htmlEl);
            if (computedStyle) {
              // 强制设置所有颜色属性为 RGB 值
              const bgColor = computedStyle.backgroundColor;
              if (bgColor && bgColor !== 'rgba(0, 0, 0, 0)' && bgColor !== 'transparent') {
                if (bgColor.startsWith('rgb')) {
                  htmlEl.style.setProperty('background-color', bgColor, 'important');
                }
              }
              
              const textColor = computedStyle.color;
              if (textColor && textColor !== 'rgba(0, 0, 0, 0)') {
                if (textColor.startsWith('rgb')) {
                  htmlEl.style.setProperty('color', textColor, 'important');
                }
              }
              
              const borderColor = computedStyle.borderColor;
              if (borderColor && borderColor !== 'rgba(0, 0, 0, 0)' && computedStyle.borderWidth !== '0px') {
                if (borderColor.startsWith('rgb')) {
                  htmlEl.style.setProperty('border-color', borderColor, 'important');
                }
              }
              
              // 处理 box-shadow（移除 oklch，使用计算后的值）
              const boxShadow = computedStyle.boxShadow;
              if (boxShadow && boxShadow !== 'none' && !boxShadow.includes('oklch')) {
                htmlEl.style.setProperty('box-shadow', boxShadow, 'important');
              }
            }
          } catch (e) {
            // 忽略无法处理的元素
          }
        });

        // 第三步：确保所有图片都已加载
        const images = clonedElement.getElementsByTagName('img');
        Array.from(images).forEach((img) => {
          if (!img.complete) {
            img.style.display = 'none';
          }
        });
      },
    });

    // 转换为 blob 并下载
    return new Promise((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (!blob) {
          reject(new Error('生成图片失败'));
          return;
        }

        // 创建下载链接
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `${filename}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        
        // 延迟释放URL，确保下载完成
        setTimeout(() => {
          URL.revokeObjectURL(url);
        }, 100);
        
        resolve();
      }, 'image/png', 1.0); // 最高质量
    });
  } catch (error) {
    console.error('导出卡片失败:', error);
    throw error;
  }
}

