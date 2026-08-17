// Confirmación mostrada antes de eliminar. El backend solamente se llama
// después de que el usuario confirme.
export default function ModalConfirmacion({ mensaje, onConfirmar, onCancelar }) {
  return (
    <div className="fondo-modal" role="dialog" aria-modal="true">
      <div className="modal">
        <p>{mensaje}</p>
        <div className="acciones-formulario">
          <button type="button" className="boton boton-peligro" onClick={onConfirmar}>
            Eliminar
          </button>
          <button type="button" className="boton" onClick={onCancelar}>
            Cancelar
          </button>
        </div>
      </div>
    </div>
  )
}
