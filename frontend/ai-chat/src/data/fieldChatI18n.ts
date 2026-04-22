export interface FieldChatStrings {
  placeholder: string
  footerHint: string
  ariaSend: string
  fieldFilledPrefix: string
  errorProcessing: string
  cancelled: string
  connectionError: string
  acceptClose: string
}

export const fieldChatI18n: Record<string, FieldChatStrings> = {
  en: {
    placeholder: 'Type your answer...',
    footerHint: 'Enter to send · Shift+Enter new line · Esc to close',
    ariaSend: 'Send',
    fieldFilledPrefix: '✓ Field filled:',
    errorProcessing: 'Could not process your message. Try again.',
    cancelled: '(cancelled)',
    connectionError: 'Connection error.',
    acceptClose: 'Accept & close',
  },
  pt: {
    placeholder: 'Digite sua resposta...',
    footerHint: 'Enter para enviar · Shift+Enter nova linha · Esc para fechar',
    ariaSend: 'Enviar',
    fieldFilledPrefix: '✓ Campo preenchido:',
    errorProcessing: 'Não consegui processar. Tente novamente.',
    cancelled: '(cancelado)',
    connectionError: 'Erro de conexão.',
    acceptClose: 'Aceitar e fechar',
  },
  es: {
    placeholder: 'Escribe tu respuesta...',
    footerHint: 'Enter para enviar · Shift+Enter nueva línea · Esc para cerrar',
    ariaSend: 'Enviar',
    fieldFilledPrefix: '✓ Campo completado:',
    errorProcessing: 'No se pudo procesar. Inténtalo de nuevo.',
    cancelled: '(cancelado)',
    connectionError: 'Error de conexión.',
    acceptClose: 'Aceptar y cerrar',
  },
  fr: {
    placeholder: 'Tapez votre réponse...',
    footerHint: 'Entrée pour envoyer · Maj+Entrée nouvelle ligne · Échap pour fermer',
    ariaSend: 'Envoyer',
    fieldFilledPrefix: '✓ Champ rempli :',
    errorProcessing: 'Impossible de traiter. Réessayez.',
    cancelled: '(annulé)',
    connectionError: 'Erreur de connexion.',
    acceptClose: 'Accepter et fermer',
  },
}
