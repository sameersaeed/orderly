import React, { useState, useEffect } from 'react';
import ItemModal from './Modal';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {  faEdit, faShoppingCart, faTrashAlt } from '@fortawesome/free-solid-svg-icons';
import axios from 'axios';
import * as toastr from 'toastr';

interface Item {
  id: string;
  name: string;
  description: string;
  price: number;
}

const Shop: React.FC<{ user: any }> = ({ user }) => {
  const [items, setItems] = useState<Item[]>([]);
  const [cart, setCart] = useState<Item[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState<Item>();

  useEffect(() => {
    switch (sessionStorage.getItem('itemActionPerformed')) {
      case 'edit':
        toastr.success('Item was updated successfully!');
        sessionStorage.removeItem('itemActionPerformed');
        break;
      case 'create':
        toastr.success('Item was created successfully!');
        sessionStorage.removeItem('itemActionPerformed');
        break;
      default:
        break;
    }
    
    const fetchItems = async () => {
      try {
        const response = await axios.get<Item[]>(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/getItems`, {
          headers: { Authorization: `Bearer ${user.Token}` }
        });
        setItems(response.data);
      } 
      catch (error) {
        console.error('ERROR: could not fetch items:', error);
      }
    };

    if (user.Token) {
      fetchItems();
    }
  }, [user.Token]);

  const addToCart = (item: Item) => {
    setCart([...cart, item]);
  };

  const removeCartItem = async (id: string) => {
    try {
      await axios.delete(`/items/${id}`, {
        headers: { Authorization: `Bearer ${user.Token}` }
      });
      setItems(items.filter(item => item.id !== id));
    } 
    catch (error) {
      console.error('ERROR: could not remove item from cart:', error);
    }
  };

  const handleSubmit = async () => {
    try {
      const orderDetails = cart.map(item => `${item.name}: ${item.price}`).join('\n');
      const message = `Order details:\n${orderDetails}`;
      await axios.post('/send', {
        name: user.name,
        email: user.Email,
        message,
        cart
      });

      toastr.success('Your order has successfully been sent for processing!');
    } 
    catch (error) {
      console.error('ERROR: could not send order:', error);
      toastr.error('ERROR: There was an issue submitting your order, please try again');
    }
  };

  const handleAddItem = async () => {
    const newItem: Item = {
      id: '',
      name: 'Item name',
      description: 'Enter a short item description here',
      price: 0.00,
    };
  
    try {
      const response = await axios.post<Item>(
        `${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/createItem`, 
        newItem,
        {
          headers: { Authorization: `Bearer ${user.Token}` }
        }
      );

      // add to items array if other items exist, otherwise create a new array
      (Array.isArray(items)) 
        ? setItems([...items, response.data]) 
        : setItems([response.data]);

        window.location.reload();
        sessionStorage.setItem('itemActionPerformed', 'create');
    } 
    catch (error) {
      console.error('ERROR: could not add new item:', error);
      toastr.error('ERROR: There was an issue with creating the item, please try again');
    }
  };

  const editItem = (item: Item) => {
    setEditingItem({ ...item });
    setShowModal(true);
  };

  const handleEditItem = async (item: Item) => {
    if (item.id) {
      try {
        const response = await axios.put(
          `${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/editItem`,
          item,
          {
            headers: { Authorization: `Bearer ${user.Token}` },
          }
        );

        setItems(items.map(item => item.id === editingItem?.id ? item : item));
        setShowModal(false);

        window.location.reload();
        sessionStorage.setItem('itemActionPerformed', 'edit');
      } 
      catch (error) {
        console.error('ERROR: could not edit item:', error);
        toastr.error('ERROR: There was an issue with editing the item, please try again');
      }
    }
    else {
      toastr.error('ERROR: Could not edit item, no item ID was found');
    }
  };

  const handleDeleteItem = async (itemID: string) => {
    try {
      await axios.delete(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/deleteItem`, {
        data: { id: itemID },
        headers: { Authorization: `Bearer ${user.Token}` }
      });

      setItems(items.filter(item => item.id !== itemID));
      toastr.success('Item deleted successfully!');
    } 
    catch (error) {
      console.error('ERROR: could not delete item:', error);
      toastr.error('ERROR: There was an issue with deleting the item, please try again');
    }
  };
  
  return (
    <div className="container mt-5">
      <h2>Items</h2>
      <div className="row">
        {items?.map((item: Item) => (
          <div key={item.id} className="col-md-4">
            <div className="card mb-4 bg-dark text-light">
              {user.IsAdmin && (
                <div className="position-absolute top-0 end-0">
                  <button className="btn btn-danger btn-sm me-2" onClick={() => handleDeleteItem(item.id)}>
                    <FontAwesomeIcon icon={faTrashAlt} />
                  </button>
                  <button className="btn btn-info btn-sm" onClick={() => editItem(item)}>
                    <FontAwesomeIcon icon={faEdit} />
                  </button>
                </div>
              )}
              <div className="card-body">
                <h5 className="card-title">{item.name}</h5>
                <p className="card-text">{item.description}</p>
                <p className="card-text">Price: ${item.price}</p>
                <button className="btn btn-primary" onClick={() => addToCart(item)}>
                  <FontAwesomeIcon icon={faShoppingCart} /> Add to Cart
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
      {user.IsAdmin && (
        <button className="btn btn-success mt-3" onClick={() => setShowModal(true)}>
          Create a new item
        </button>
      )}
      {showModal && <ItemModal onClose={() => setShowModal(false)} onSave={editingItem ? handleEditItem : handleAddItem} item={editingItem || { id: '', name: '', description: '', price: 0 }} />}
      {cart.length > 0 && (
        <div className="mt-4">
          <h3>Cart</h3>
          <ul className="list-group">
            {cart.map(item => (
              <li key={item.id} className="list-group-item d-flex justify-content-between align-items-center">
                {item.name}
                <button className="btn btn-danger btn-sm" onClick={() => removeCartItem(item.id)}>
                  Remove
                </button>
              </li>
            ))}
          </ul>
          <button className="btn btn-success mt-3" onClick={handleSubmit}>
            Submit Order
          </button>
        </div>
      )}
    </div>
  );
}

export default Shop;
