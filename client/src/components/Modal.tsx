import React, { useState } from 'react';
import { Types } from 'mongoose';

interface Item {
  _id: string;
  Name: string;
  Description: string;
  Price: number;
}

interface ModalProps {
  onClose: () => void;
  onSave: (item: Item) => Promise<void>;
  item?: Item;
}

const Modal: React.FC<ModalProps> = ({ onClose, onSave, item }) => {
  const [Name, setName] = useState(item?.Name || '');
  const [Description, setDescription] = useState(item?.Description || '');
  const [Price, setPrice] = useState(item?.Price || 0);

  const handleSave = () => {
    const newItem: Item = {
      _id: '',
      Name,
      Description,
      Price
    };
    onSave(newItem);
  };

  return (
    <div className="modal show fade" style={{ display: 'block' }} role="dialog" aria-labelledby="modalTitle" aria-hidden="true">
      <div className="modal-dialog modal-dialog-centered" role="document">
        <div className="modal-content bg-dark text-light">
          <div className="modal-header">
            <h5 className="modal-title" id="modalTitle">{item?._id !== '' ? 'Edit Item' : 'Add New Item'}</h5>
            <button type="button" className="btn-close btn-close-white" onClick={onClose}></button>
          </div>
          <div className="modal-body">
            <form>
              <div className="form-group row">
                <label htmlFor="itemName" className="col-sm-3 col-form-label">Name</label>
                <div className="col-sm-9">
                  <input type="text" className="form-control" id="itemName" value={Name} onChange={(e) => setName(e.target.value)} />
                </div>
              </div>
              <div className="form-group row">
                <label htmlFor="itemDescription" className="col-sm-3 col-form-label">Description</label>
                <div className="col-sm-9">
                  <input type="text" className="form-control" id="itemDescription" value={Description} onChange={(e) => setDescription(e.target.value)} />
                </div>
              </div>
              <div className="form-group row">
                <label htmlFor="itemPrice" className="col-sm-3 col-form-label">Price</label>
                <div className="col-sm-9">
                  <input type="number" className="form-control" id="itemPrice" value={Price} onChange={(e) => setPrice(Number(e.target.value))} />
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
