import React, { useState } from 'react';

interface Item {
  id: string;
  name: string;
  description: string;
  price: number;
}

interface ModalProps {
  onClose: () => void;
  onSave: (item: Item) => Promise<void>;
  item: Item;
}

const Modal: React.FC<ModalProps> = ({ onClose, onSave, item }) => {
  const [name, setName] = useState(item.name || '');
  const [description, setDescription] = useState(item.description || '');
  const [price, setPrice] = useState(item.price || 0);

  const handleSave = () => {
    const newItem: Item = {
      id: item.id,
      name,
      description,
      price,
    };
    onSave(newItem);
  };

  return (
    <div className="modal show fade" style={{ display: 'block' }} role="dialog" aria-labelledby="modalTitle" aria-hidden="true">
      <div className="modal-dialog modal-dialog-centered" role="document">
        <div className="modal-content bg-dark text-light">
          <div className="modal-header">
            <h5 className="modal-title" id="modalTitle">{item.id ? 'Edit Item' : 'Add New Item'}</h5>
            <button type="button" className="btn-close btn-close-white" onClick={onClose}></button>
          </div>
          <div className="modal-body">
            <form>
              <div className="form-group row">
                <label htmlFor="itemName" className="col-sm-3 col-form-label">Name</label>
                <div className="col-sm-9">
                  <input type="text" className="form-control" id="itemName" value={name} onChange={(e) => setName(e.target.value)} />
                </div>
              </div>
              <div className="form-group row">
                <label htmlFor="itemDescription" className="col-sm-3 col-form-label">Description</label>
                <div className="col-sm-9">
                  <input type="text" className="form-control" id="itemDescription" value={description} onChange={(e) => setDescription(e.target.value)} />
                </div>
              </div>
              <div className="form-group row">
                <label htmlFor="itemPrice" className="col-sm-3 col-form-label">Price</label>
                <div className="col-sm-9">
                  <input type="number" className="form-control" id="itemPrice" value={price} onChange={(e) => setPrice(Number(e.target.value))} />
                </div>
              </div>
            </form>
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={onClose}>Close</button>
            <button type="button" className="btn btn-primary" onClick={handleSave}>Save</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Modal;