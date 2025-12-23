/**
 * 中文翻译文件
 */

export default {
  // 通用
  common: {
    confirm: '确认',
    cancel: '取消',
    save: '保存',
    delete: '删除',
    edit: '编辑',
    back: '返回',
    next: '下一步',
    loading: '加载中...',
    error: '错误',
    success: '成功',
    report: '学习报告',
  },

  // Header
  header: {
    title: '小探号',
    favorites: '我的收藏',
  },

  // 首页
  home: {
    title: 'TanGo 探索',
    photoButton: '拍照探索',
    voiceButton: '语音输入',
    quickCapture: '快速拍照',
    cardScience: '科学认知',
    cardHumanities: '人文素养',
    cardLanguage: '语言能力',
    littleStarMessage: '拍一拍，发现有趣的知识吧～',
  },

  // 拍照页面
  capture: {
    title: '拍照探索',
    takePhoto: '拍照',
    selectFromAlbum: '从相册选择',
    voiceInput: '语音输入',
    processing: '识别中...',
    aiAutoDetect: 'AI自动识别',
    identifyError: '识别失败，请重试',
    identifyErrorDetail: '识别失败: {{error}}\n\n请检查：\n1. 图片是否清晰\n2. 网络连接是否正常\n3. 稍后重试',
    voiceNotSupported: '您的浏览器不支持语音识别功能',
  },

  // 结果页面
  result: {
    title: '探索结果',
    identifiedAs: '识别为',
    cards: '知识卡片',
    continueChat: '继续对话',
    collect: '收藏',
    collected: '已收藏',
    foundNewFriend: '你发现了一个新朋友！',
    itsA: '这是一个',
    aiCompanionSays: 'AI小伙伴说：',
    aiCompanionMessage: '哇！这是一个{{objectName}}！让我们探索它的秘密吧！',
  },

  // 收藏页面
  collection: {
    title: '我的收藏',
    subtitle: '继续探索你的收藏吧！',
    empty: '还没有收藏任何卡片',
    export: '导出',
    clearAll: '清空所有',
    exportAll: '导出全部',
    exportAllTitle: '导出所有卡片',
    parentMode: '家长模式',
    clearAllHint: '仅在家长模式下可用',
    littleStarSays: '小探星说：',
    littleStarMessage: '去探索有趣的知识，收藏更多喜欢的卡片吧！我在等待你的发现！✨',
    exportError: '导出失败，请重试',
    exportCardTitle: '导出卡片',
    emptyMessage: '还没有收藏任何卡片，快去探索吧！',
    reExplore: '重新探索',
    collect: '收藏',
    uncollect: '取消收藏',
    category: {
      all: '全部',
      natural: '自然类',
      life: '生活类',
      humanities: '人文类',
      science: '科学认知',
      poetry: '古诗词/人文',
      english: '英语表达',
    },
  },

  // 设置页面
  settings: {
    title: '设置',
    language: '语言',
    languageDesc: '选择应用显示语言',
    grade: '年级',
    gradeDesc: '选择孩子当前年级',
    about: '关于',
    version: '版本',
    appDescription: 'TanGo - 探索世界的知识卡片应用',
    gradeK1: '幼儿园小班',
    gradeK2: '幼儿园中班',
    gradeK3: '幼儿园大班',
    gradeG1: '一年级',
    gradeG2: '二年级',
    gradeG3: '三年级',
    gradeG4: '四年级',
    gradeG5: '五年级',
    gradeG6: '六年级',
    gradeG7: '七年级',
    gradeG8: '八年级',
    gradeG9: '九年级',
    gradeG10: '十年级',
    gradeG11: '十一年级',
    gradeG12: '十二年级',
  },

  // 对话
  conversation: {
    placeholder: '输入消息...',
    send: '发送',
    voiceInput: '语音输入',
    imageInput: '图片输入',
    thinking: '思考中...',
    error: '发送失败，请重试',
    generatingCards: '正在为您生成知识卡片...',
    generatingCardsWait: '正在生成知识卡片，请稍候...',
    generateCardsError: '生成卡片失败：{{error}}。您可以稍后通过对话重新生成。',
    generateAnswerError: '抱歉，生成回答时出现错误：{{error}}',
    sendMessageError: '抱歉，发送消息失败：{{error}}。请检查网络连接后重试。',
    sendVoiceError: '抱歉，发送语音消息失败：{{error}}。请重试。',
    sendImageError: '抱歉，发送图片失败：{{error}}。请检查图片格式和大小后重试。',
    unknownError: '未知错误',
    voiceError: '语音识别失败',
    voiceNetworkError: '网络连接失败，请检查网络后重试',
    voiceNoSpeech: '未检测到语音，请重试',
    voiceAudioError: '无法访问麦克风，请检查权限设置',
    voiceNotAllowed: '麦克风权限被拒绝，请在浏览器设置中允许访问',
    voiceStartError: '启动语音识别失败，请重试',
    voiceNotSupported: '您的浏览器不支持语音识别功能',
    voiceTimeout: '语音识别超时，请重试',
    voiceStop: '点击停止语音识别',
  },

  // 卡片类型
  cardType: {
    science: '科学认知卡',
    poetry: '古诗词/人文卡',
    english: '英语表达卡',
  },

  // 分享页面
  share: {
    invalidLink: '分享链接无效',
    loadError: '加载分享数据失败',
    loadFailed: '加载失败',
    title: '孩子的探索成果',
    createdAt: '创建时间',
    expiresAt: '过期时间',
    noRecords: '暂无探索记录',
  },

  // Little Star
  littleStar: {
    name: '小探星',
  },

  // 报告页面
  report: {
    weeklyReport: '周报',
    greeting: '你好，小探号！',
    subtitle: '你做得很好！看看你这周的成长吧。',
    explorationStars: '探索次数',
    keepExploring: '继续探索！',
    totalFavorites: '收藏总数',
    greatCollection: '收藏很棒！',
    littleExpert: '小小专家',
    natureMaster: '自然大师',
    levelUp: '升级了！🚀',
    knowledgeMap: '知识地图',
    total: '总数',
    categoryNatural: '自然类',
    categoryLife: '生活类',
    categoryHumanities: '人文类',
    items: '项',
    recentFavorites: '最近收藏',
    recentFavoritesMessage: '最近收藏了 {{totalCollectedCards}} 张卡片',
    noCards: '还没有收藏任何卡片',
  },
};
