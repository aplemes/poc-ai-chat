export interface ChatStrings {
  code: string
  placeholder: string
  greeting: string
  formFilled: string
  errorProcessing: string
  cancelled: string
  errorConnection: string
  footerHint: string
  confirmTitle: string
  confirmBtn: string
  correctBtn: string
  correctionHint: string
  ariaOpenChat: string
  ariaOpenChatNew: string
  ariaCloseChat: string
  ariaResetChat: string
  ariaConversation: string
  ariaSend: string
}

export const chatI18n: ChatStrings[] = [
  {
    code: 'pt',
    placeholder: 'Escreva sua mensagem...',
    greeting: 'Olá! Como posso ajudar você hoje?',
    formFilled: 'Formulário preenchido! Revise os campos.',
    errorProcessing: 'Não consegui processar. Tente novamente.',
    cancelled: '(cancelado)',
    errorConnection: 'Sem conexão com o servidor.',
    footerHint: 'Enter para enviar · Shift+Enter nova linha',
    confirmTitle: 'Revisar antes de preencher',
    confirmBtn: 'Confirmar e preencher',
    correctBtn: 'Corrigir pelo chat',
    correctionHint: 'Precisa ajustar algo? É só me dizer pelo chat.',
    ariaOpenChat: 'Abrir assistente',
    ariaOpenChatNew: 'Abrir assistente – nova mensagem',
    ariaCloseChat: 'Fechar chat',
    ariaResetChat: 'Nova conversa',
    ariaConversation: 'Conversa',
    ariaSend: 'Enviar',
  },
  {
    code: 'en',
    placeholder: 'Write your message...',
    greeting: 'Hello! How can I help you today?',
    formFilled: 'Form filled! Please review the fields.',
    errorProcessing: 'Could not process your message. Try again.',
    cancelled: '(cancelled)',
    errorConnection: 'Could not connect to the server.',
    footerHint: 'Enter to send · Shift+Enter for new line',
    confirmTitle: 'Review before filling',
    confirmBtn: 'Confirm & fill',
    correctBtn: 'Correct via chat',
    correctionHint: 'Need to adjust something? Just tell me in the chat.',
    ariaOpenChat: 'Open assistant',
    ariaOpenChatNew: 'Open assistant – new message',
    ariaCloseChat: 'Close chat',
    ariaResetChat: 'New conversation',
    ariaConversation: 'Conversation',
    ariaSend: 'Send',
  },
  {
    code: 'es',
    placeholder: 'Escribe tu mensaje...',
    greeting: '¡Hola! ¿Cómo puedo ayudarte hoy?',
    formFilled: '¡Formulario completado! Revisa los campos.',
    errorProcessing: 'No se pudo procesar. Inténtalo de nuevo.',
    cancelled: '(cancelado)',
    errorConnection: 'No se pudo conectar al servidor.',
    footerHint: 'Enter para enviar · Shift+Enter nueva línea',
    confirmTitle: 'Revisar antes de completar',
    confirmBtn: 'Confirmar y completar',
    correctBtn: 'Corregir por chat',
    correctionHint: '¿Necesitas ajustar algo? Solo dímelo por el chat.',
    ariaOpenChat: 'Abrir asistente',
    ariaOpenChatNew: 'Abrir asistente – nuevo mensaje',
    ariaCloseChat: 'Cerrar chat',
    ariaResetChat: 'Nueva conversación',
    ariaConversation: 'Conversación',
    ariaSend: 'Enviar',
  },
  {
    code: 'fr',
    placeholder: 'Écrivez votre message...',
    greeting: 'Bonjour ! Comment puis-je vous aider ?',
    formFilled: 'Formulaire rempli ! Vérifiez les champs.',
    errorProcessing: 'Impossible de traiter. Réessayez.',
    cancelled: '(annulé)',
    errorConnection: 'Connexion au serveur impossible.',
    footerHint: 'Entrée pour envoyer · Maj+Entrée nouvelle ligne',
    confirmTitle: 'Vérifier avant de remplir',
    confirmBtn: 'Confirmer et remplir',
    correctBtn: 'Corriger via chat',
    correctionHint: "Besoin d'ajuster ? Dites-le moi dans le chat.",
    ariaOpenChat: "Ouvrir l'assistant",
    ariaOpenChatNew: "Ouvrir l'assistant – nouveau message",
    ariaCloseChat: 'Fermer le chat',
    ariaResetChat: 'Nouvelle conversation',
    ariaConversation: 'Conversation',
    ariaSend: 'Envoyer',
  },
]

export const langNames: Record<string, string> = {
  pt: 'Português',
  en: 'English',
  es: 'Español',
  fr: 'Français',
}
