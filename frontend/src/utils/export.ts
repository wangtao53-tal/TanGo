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
    // 使用 html2canvas 捕获元素，优化配置确保清晰度
    const canvas = await html2canvas(element, {
      backgroundColor: '#ffffff',
      scale: 3, // 提高图片清晰度（从2提升到3）
      useCORS: true,
      allowTaint: false,
      logging: false,
      width: element.scrollWidth,
      height: element.scrollHeight,
      windowWidth: element.scrollWidth,
      windowHeight: element.scrollHeight,
      // 确保所有样式都被渲染
      onclone: (clonedDoc) => {
        const clonedElement = clonedDoc.getElementById(elementId);
        if (!clonedElement || !clonedDoc.defaultView) return;

        // 第一步：移除所有包含 oklch 的样式表规则
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
                        try {
                          if (nestedRule.cssText && nestedRule.cssText.includes('oklch')) {
                            if (rule instanceof CSSMediaRule) {
                              rule.deleteRule(j);
                            } else if (rule instanceof CSSKeyframesRule) {
                              rule.deleteRule(j);
                            }
                          }
                        } catch (e) {
                          // 忽略无法删除的嵌套规则
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
            }
          });
        } catch (e) {
          // 忽略样式表处理错误
        }

        // 第二步：遍历所有元素，清理 style 属性中的 oklch，并使用 getComputedStyle 设置 RGB 值
        const allElements = clonedElement.querySelectorAll('*');
        const elementsArray = [clonedElement, ...Array.from(allElements)];
        
        if (!clonedDoc.defaultView) return;
        
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
            if (!clonedDoc.defaultView) return;
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
