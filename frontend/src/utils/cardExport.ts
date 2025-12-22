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
    scale = 4, // 提高清晰度到4倍，使导出的图片更清晰
    backgroundColor = '#ffffff',
    filename = `card-${Date.now()}`,
  } = options;

  const element = document.getElementById(elementId);
  if (!element) {
    throw new Error(`找不到ID为 ${elementId} 的元素`);
  }

  try {
    // 使用 html2canvas 捕获元素，优化配置确保清晰度
    // 使用scrollWidth和scrollHeight确保包含所有内容
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
      scrollX: 0,
      scrollY: 0,
      // 确保完整捕获所有内容，包括底部
      removeContainer: false,
      // 忽略某些元素（已在onclone中处理）
      ignoreElements: (element) => {
        // 忽略所有按钮
        if (element.tagName === 'BUTTON') {
          return true;
        }
        // 忽略包含特定图标的元素
        const hasIcon = element.querySelector && element.querySelector('span[class*="material-symbols"]');
        if (hasIcon && (element.classList.contains('absolute') || element.classList.contains('fixed'))) {
          return true;
        }
        return false;
      },
      // 确保所有样式都被渲染
      onclone: (clonedDoc) => {
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
              console.warn('无法访问样式表:', e);
            }
          });
        } catch (e) {
          console.warn('处理样式表时出错:', e);
        }

        const clonedElement = clonedDoc.getElementById(elementId);
        if (!clonedElement || !clonedDoc.defaultView) return;

        // 隐藏滚动条，避免遮挡内容
        // 使用CSS样式注入来隐藏所有滚动条，但不隐藏内容
        const style = clonedDoc.createElement('style');
        style.textContent = `
          * {
            scrollbar-width: none !important;
            -ms-overflow-style: none !important;
          }
          *::-webkit-scrollbar {
            display: none !important;
            width: 0 !important;
            height: 0 !important;
            background: transparent !important;
          }
          *::-webkit-scrollbar-track {
            display: none !important;
          }
          *::-webkit-scrollbar-thumb {
            display: none !important;
          }
        `;
        clonedDoc.head.appendChild(style);

        // 只隐藏滚动条，但保持内容可见（使用overflow: visible或auto）
        clonedElement.style.setProperty('overflow', 'visible', 'important');
        clonedElement.style.setProperty('overflow-x', 'visible', 'important');
        clonedElement.style.setProperty('overflow-y', 'visible', 'important');
        
        // 确保所有子元素内容可见，只隐藏滚动条
        const allChildren = clonedElement.querySelectorAll('*');
        allChildren.forEach((child) => {
          const htmlChild = child as HTMLElement;
          // 只设置滚动条隐藏，不设置overflow hidden，避免内容被裁剪
          htmlChild.style.setProperty('scrollbar-width', 'none', 'important');
          htmlChild.style.setProperty('-ms-overflow-style', 'none', 'important');
          // 如果元素有overflow hidden，改为visible或auto
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlChild);
          if (computedStyle && computedStyle.overflow === 'hidden') {
            htmlChild.style.setProperty('overflow', 'visible', 'important');
          }
        });

        // 隐藏所有按钮和UI元素，避免遮挡内容
        // 1. 隐藏所有按钮元素
        const buttons = clonedElement.querySelectorAll('button');
        buttons.forEach((btn) => {
          (btn as HTMLElement).style.display = 'none';
        });

        // 2. 隐藏所有包含 material-symbols 图标的元素（通常是按钮图标）
        const iconElements = clonedElement.querySelectorAll('span[class*="material-symbols"], [class*="material-symbols"]');
        iconElements.forEach((icon) => {
          const parent = icon.parentElement;
          if (parent && (parent.tagName === 'BUTTON' || parent.classList.contains('absolute') || parent.classList.contains('fixed'))) {
            (parent as HTMLElement).style.display = 'none';
          }
        });

        // 3. 隐藏导航按钮（chevron_left, chevron_right）
        const navButtons = clonedElement.querySelectorAll('[class*="chevron"], [aria-label*="上一张"], [aria-label*="下一张"], [aria-label*="切换"]');
        navButtons.forEach((btn) => {
          (btn as HTMLElement).style.display = 'none';
        });

        // 4. 隐藏导出按钮
        const exportButtons = clonedElement.querySelectorAll('[aria-label*="导出"], [title*="导出"], [class*="download"]');
        exportButtons.forEach((btn) => {
          (btn as HTMLElement).style.display = 'none';
        });

        // 5. 隐藏指示器（分页点）
        const indicators = clonedElement.querySelectorAll('[aria-label*="切换"], [class*="indicator"], [class*="rounded-full"][class*="bg-"][class*="w-"]');
        indicators.forEach((indicator) => {
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(indicator as HTMLElement);
          if (computedStyle && (computedStyle.width === '8px' || computedStyle.width === '0.5rem' || computedStyle.width === '24px' || computedStyle.width === '1.5rem')) {
            (indicator as HTMLElement).style.display = 'none';
          }
        });

        // 6. 隐藏所有绝对定位的元素（通常在右上角、边缘或覆盖在内容上）
        const allAbsoluteElements = clonedElement.querySelectorAll('*');
        allAbsoluteElements.forEach((el) => {
          const htmlEl = el as HTMLElement;
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlEl);
          if (computedStyle && (computedStyle.position === 'absolute' || computedStyle.position === 'fixed')) {
            // 隐藏所有绝对定位的按钮、图标或交互元素
            if (htmlEl.tagName === 'BUTTON' || 
                htmlEl.querySelector('span[class*="material-symbols"]') || 
                htmlEl.classList.contains('absolute') ||
                htmlEl.classList.contains('fixed')) {
              htmlEl.style.display = 'none';
            }
          }
        });

        // 7. 隐藏操作按钮区域（底部包含"听"、"收藏"等按钮的区域）
        const actionAreas = clonedElement.querySelectorAll('[class*="flex"][class*="justify-between"][class*="items-center"]');
        actionAreas.forEach((area) => {
          const htmlArea = area as HTMLElement;
          // 检查是否包含按钮
          if (htmlArea.querySelector('button')) {
            htmlArea.style.display = 'none';
          }
        });

        // 优化卡片样式，使其在导出时更美观
        // 增加卡片整体的内边距，让文字与边框有更多空间
        const cardElement = clonedElement;
        if (cardElement) {
          // 为卡片添加额外的内边距（在原有基础上增加）
          const currentPadding = clonedDoc.defaultView?.getComputedStyle(cardElement).padding || '0px';
          const paddingValue = parseInt(currentPadding) || 0;
          // 增加至少32px的内边距（如果原来没有或很小）
          const newPadding = Math.max(paddingValue, 32);
          cardElement.style.setProperty('padding', `${newPadding}px`, 'important');
        }

        // 增加内容区域的内边距，确保文字不贴边
        const cardContent = clonedElement.querySelector('[class*="p-6"], [class*="p-"], [class*="px-"], [class*="py-"]');
        if (cardContent) {
          const htmlContent = cardContent as HTMLElement;
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlContent);
          if (computedStyle) {
            // 获取当前padding值
            const currentPadding = computedStyle.padding || '24px';
            const paddingMatch = currentPadding.match(/(\d+)px/);
            const paddingValue = paddingMatch ? parseInt(paddingMatch[1]) : 24;
            // 增加内边距到至少48px（原来p-6是24px，增加到48px，确保文字与边框有足够距离）
            const newPadding = Math.max(paddingValue, 48);
            htmlContent.style.setProperty('padding', `${newPadding}px`, 'important');
            // 特别增加底部内边距，确保底部内容完整显示
            htmlContent.style.setProperty('padding-bottom', '60px', 'important');
          } else {
            // 如果没有计算样式，直接设置较大的内边距
            htmlContent.style.setProperty('padding', '48px', 'important');
            htmlContent.style.setProperty('padding-bottom', '60px', 'important');
          }
        }

        // 为所有文本容器增加内边距
        const textContainers = clonedElement.querySelectorAll('p, h1, h2, h3, h4, h5, h6, div[class*="bg-"], div[class*="border"]');
        textContainers.forEach((container) => {
          const htmlContainer = container as HTMLElement;
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlContainer);
          if (computedStyle) {
            const currentPadding = computedStyle.padding || '0px';
            const paddingMatch = currentPadding.match(/(\d+)px/);
            const paddingValue = paddingMatch ? parseInt(paddingMatch[1]) : 0;
            // 如果内边距小于16px，增加到16px
            if (paddingValue < 16 && htmlContainer.textContent && htmlContainer.textContent.trim().length > 0) {
              htmlContainer.style.setProperty('padding', '16px', 'important');
            }
          }
        });

        // 确保卡片有足够的圆角和阴影，使其更美观
        if (cardElement) {
          cardElement.style.setProperty('border-radius', '2.5rem', 'important');
          cardElement.style.setProperty('box-shadow', '0 10px 40px rgba(0, 0, 0, 0.1)', 'important');
          cardElement.style.setProperty('overflow', 'visible', 'important'); // 改为visible，确保内容不被裁剪
          // 确保背景色为白色
          cardElement.style.setProperty('background-color', '#ffffff', 'important');
          // 确保卡片高度自适应内容，不被压缩
          cardElement.style.setProperty('min-height', 'auto', 'important');
          cardElement.style.setProperty('height', 'auto', 'important');
          // 确保宽度固定，避免布局问题
          if (cardElement.style.width) {
            cardElement.style.setProperty('width', cardElement.style.width, 'important');
          }
          
          // 获取卡片的边框颜色（从现有边框中提取）
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(cardElement);
          let borderColor = '#000000'; // 默认黑色
          if (computedStyle) {
            borderColor = computedStyle.borderColor || computedStyle.borderTopColor || '#000000';
            // 如果边框颜色是rgba(0,0,0,0)或transparent，使用默认颜色
            if (borderColor === 'rgba(0, 0, 0, 0)' || borderColor === 'transparent') {
              // 尝试从类名中推断颜色
              if (cardElement.className.includes('border-science-green')) {
                borderColor = '#10B981'; // 绿色
              } else if (cardElement.className.includes('border-sunny-orange')) {
                borderColor = '#F97316'; // 橙色
              } else if (cardElement.className.includes('border-sky-blue')) {
                borderColor = '#0EA5E9'; // 蓝色
              } else {
                borderColor = '#000000'; // 默认黑色
              }
            }
          }
          
          // 确保底部边框可见且足够粗
          const currentBorderBottom = computedStyle?.borderBottomWidth || '0px';
          const borderBottomValue = parseInt(currentBorderBottom) || 0;
          if (borderBottomValue < 4) {
            cardElement.style.setProperty('border-bottom-width', '4px', 'important');
            cardElement.style.setProperty('border-bottom-style', 'solid', 'important');
            cardElement.style.setProperty('border-bottom-color', borderColor, 'important');
          }
          
          // 添加底部内边距，确保底部内容完整显示且有边框包裹
          const currentPaddingBottom = computedStyle?.paddingBottom || '0px';
          const paddingBottomValue = parseInt(currentPaddingBottom) || 0;
          if (paddingBottomValue < 50) {
            cardElement.style.setProperty('padding-bottom', '50px', 'important');
          }
        }

        // 优化内容区域的样式和间距
        const contentAreas = clonedElement.querySelectorAll('[class*="p-"], [class*="px-"], [class*="py-"], [class*="mb-"], [class*="mt-"]');
        contentAreas.forEach((area) => {
          const htmlArea = area as HTMLElement;
          if (htmlArea) {
            // 确保文本区域有足够的对比度
            const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlArea);
            if (computedStyle) {
              // 确保背景色
              if (!computedStyle.backgroundColor || computedStyle.backgroundColor === 'rgba(0, 0, 0, 0)') {
                htmlArea.style.setProperty('background-color', '#ffffff', 'important');
              }
              
              // 增加段落和标题的上下间距，让内容更易读
              if (htmlArea.tagName === 'P' || htmlArea.tagName === 'H1' || htmlArea.tagName === 'H2' || htmlArea.tagName === 'H3' || htmlArea.tagName === 'H4') {
                const currentMargin = computedStyle.marginBottom || '0px';
                const marginValue = parseInt(currentMargin) || 0;
                if (marginValue < 16) {
                  htmlArea.style.setProperty('margin-bottom', '16px', 'important');
                }
              }
            }
          }
        });

        // 为所有文本块增加额外的内边距和间距
        const textBlocks = clonedElement.querySelectorAll('p, h1, h2, h3, h4, h5, h6, div[class*="space-y"], div[class*="mb-"]');
        textBlocks.forEach((block) => {
          const htmlBlock = block as HTMLElement;
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlBlock);
          if (computedStyle) {
            // 如果文本块没有足够的padding，添加一些
            const padding = computedStyle.padding || '0px';
            const paddingMatch = padding.match(/(\d+)px/);
            const paddingValue = paddingMatch ? parseInt(paddingMatch[1]) : 0;
            if (paddingValue < 12 && htmlBlock.textContent && htmlBlock.textContent.trim().length > 0) {
              htmlBlock.style.setProperty('padding-left', '12px', 'important');
              htmlBlock.style.setProperty('padding-right', '12px', 'important');
              htmlBlock.style.setProperty('padding-top', '8px', 'important');
              htmlBlock.style.setProperty('padding-bottom', '8px', 'important');
            }
            
            // 确保段落有足够的上下间距
            if (htmlBlock.tagName === 'P') {
              const marginBottom = computedStyle.marginBottom || '0px';
              const marginValue = parseInt(marginBottom) || 0;
              if (marginValue < 12) {
                htmlBlock.style.setProperty('margin-bottom', '12px', 'important');
              }
            }
          }
        });

        // 确保所有内容区域都可见，不被裁剪
        const contentDivs = clonedElement.querySelectorAll('div[class*="flex"], div[class*="bg-"], div[class*="p-"]');
        contentDivs.forEach((div) => {
          const htmlDiv = div as HTMLElement;
          htmlDiv.style.setProperty('overflow', 'visible', 'important');
          htmlDiv.style.setProperty('overflow-x', 'visible', 'important');
          htmlDiv.style.setProperty('overflow-y', 'visible', 'important');
          // 确保底部内容不被裁剪，添加底部内边距
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlDiv);
          if (computedStyle) {
            const currentPaddingBottom = computedStyle.paddingBottom || '0px';
            const paddingBottomValue = parseInt(currentPaddingBottom) || 0;
            if (paddingBottomValue < 20) {
              htmlDiv.style.setProperty('padding-bottom', '20px', 'important');
            }
          }
        });

        // 确保最后一个元素有足够的底部间距，并添加底部边框装饰
        const lastChild = clonedElement.lastElementChild;
        if (lastChild) {
          const htmlLastChild = lastChild as HTMLElement;
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(htmlLastChild);
          if (computedStyle) {
            const currentMarginBottom = computedStyle.marginBottom || '0px';
            const marginBottomValue = parseInt(currentMarginBottom) || 0;
            if (marginBottomValue < 50) {
              htmlLastChild.style.setProperty('margin-bottom', '50px', 'important');
            }
            const currentPaddingBottom = computedStyle.paddingBottom || '0px';
            const paddingBottomValue = parseInt(currentPaddingBottom) || 0;
            if (paddingBottomValue < 50) {
              htmlLastChild.style.setProperty('padding-bottom', '50px', 'important');
            }
          }
        }

        // 为卡片底部添加一个装饰性边框元素，确保底部有明确的边界
        // 获取边框颜色（使用之前计算的borderColor）
        let bottomBorderColor = '#000000'; // 默认黑色
        if (cardElement) {
          const computedStyle = clonedDoc.defaultView?.getComputedStyle(cardElement);
          if (computedStyle) {
            const borderTopColor = computedStyle.borderTopColor || computedStyle.borderColor;
            if (borderTopColor && borderTopColor !== 'rgba(0, 0, 0, 0)' && borderTopColor !== 'transparent') {
              bottomBorderColor = borderTopColor;
            } else {
              // 从类名推断颜色
              const className = cardElement.className || '';
              if (className.includes('border-science-green')) {
                bottomBorderColor = '#10B981'; // 绿色
              } else if (className.includes('border-sunny-orange')) {
                bottomBorderColor = '#F97316'; // 橙色
              } else if (className.includes('border-sky-blue')) {
                bottomBorderColor = '#0EA5E9'; // 蓝色
              }
            }
          }
        }
        
        // 找到内容区域（通常是第一个包含p-6或p-的div），在它内部添加底部边框
        const contentArea = clonedElement.querySelector('[class*="p-6"], [class*="p-"], [class*="px-"], [class*="py-"]');
        const targetParent = contentArea || clonedElement;
        
        // 创建底部边框装饰元素
        const bottomBorder = clonedDoc.createElement('div');
        bottomBorder.style.setProperty('width', 'calc(100% - 96px)', 'important'); // 减去左右padding
        bottomBorder.style.setProperty('height', '4px', 'important');
        bottomBorder.style.setProperty('background-color', bottomBorderColor, 'important');
        bottomBorder.style.setProperty('margin-top', '40px', 'important');
        bottomBorder.style.setProperty('margin-bottom', '20px', 'important');
        bottomBorder.style.setProperty('margin-left', 'auto', 'important');
        bottomBorder.style.setProperty('margin-right', 'auto', 'important');
        bottomBorder.style.setProperty('border-radius', '2px', 'important');
        bottomBorder.style.setProperty('flex-shrink', '0', 'important');
        bottomBorder.style.setProperty('display', 'block', 'important');
        
        // 添加到内容区域的底部（在最后一个子元素之后）
        if (targetParent && targetParent instanceof HTMLElement) {
          targetParent.appendChild(bottomBorder);
        } else if (clonedElement instanceof HTMLElement) {
          clonedElement.appendChild(bottomBorder);
        }

        // 第二步：遍历所有元素，清理 style 属性中的 oklch，并使用 getComputedStyle 设置 RGB 值
        const allElements = clonedElement.querySelectorAll('*');
        const elementsArray = [clonedElement, ...Array.from(allElements)];
        
        if (!clonedDoc.defaultView) return;
        
        elementsArray.forEach((el) => {
          const htmlEl = el as HTMLElement;
          if (!htmlEl) return;

          try {
            // 移除transform相关的hover效果，确保导出时样式稳定
            htmlEl.style.setProperty('transform', 'none', 'important');
            htmlEl.style.setProperty('transition', 'none', 'important');

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

    // 检查canvas是否有内容
    if (!canvas || canvas.width === 0 || canvas.height === 0) {
      throw new Error('生成的canvas为空，请检查元素是否可见');
    }

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
        link.style.display = 'none'; // 隐藏下载链接，避免影响页面
        document.body.appendChild(link);
        link.click();
        
        // 延迟清理，确保下载完成
        setTimeout(() => {
          document.body.removeChild(link);
          URL.revokeObjectURL(url);
          resolve();
        }, 200);
      }, 'image/png', 1.0); // 最高质量
    });
  } catch (error) {
    console.error('导出卡片失败:', error);
    throw error;
  }
}

