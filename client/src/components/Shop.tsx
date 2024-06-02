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

interface CartItem {
  item: Item;
  quantity: number;
}

const Shop: React.FC<{ user: any }> = ({ user }) => {
  const [items, setItems] = useState<Item[]>([]);
  const [cart, setCart] = useState<CartItem[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState<Item>();
  const [totalPrice, setTotalPrice] = useState(0);

  useEffect(() => {
    // updating cart and total price from local storage
    const cart = localStorage.getItem('cart');
    if (cart) {
      setCart(JSON.parse(cart));
    }

    const totalPrice = localStorage.getItem('totalPrice');
    if (totalPrice) {
      setTotalPrice(Number(totalPrice));
    }
  }
  , []);

  useEffect(() => {
    // show toastr message after item action is performed and page refreshes
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


  const saveCartInfo = (cart: CartItem[]) => {
    setCart(cart);
    const totalPrice = cart.reduce((total, cartItem) => total + cartItem.quantity * cartItem.item.price, 0);
    setTotalPrice(totalPrice);

    localStorage.setItem('cart', JSON.stringify(cart));
    localStorage.setItem('totalPrice', totalPrice.toString());
  };

  const addToCart = (item: Item) => {
    const existingItem = cart.find(cartItem => cartItem.item.id === item.id);
    if (existingItem) {
      const updatedCart = cart.map(cartItem =>
        cartItem.item.id === existingItem.item.id ? { ...cartItem, quantity: cartItem.quantity + 1 } : cartItem
      );
      saveCartInfo(updatedCart);
    } 
    else {
      const updatedCart = [...cart, { item, quantity: 1 }];
      saveCartInfo(updatedCart);
    }
  };

  const removeCartItem = (id: string) => {
    const updatedCart = cart.filter(cartItem => cartItem.item.id !== id);
    saveCartInfo(updatedCart);
  };

  const incrementCartItemQty = (id: string) => {
    const updatedCart = cart.map(cartItem =>
      cartItem.item.id === id ? { ...cartItem, quantity: cartItem.quantity + 1 } : cartItem
    );
    saveCartInfo(updatedCart);
  };

  const decrementCartItemQty = (id: string) => {
    const existingItem = cart.find(cartItem => cartItem.item.id === id);
    if (existingItem && existingItem.quantity > 1) {
      const updatedCart = cart.map(cartItem =>
        cartItem.item.id === id ? { ...cartItem, quantity: cartItem.quantity - 1 } : cartItem
      );
      saveCartInfo(updatedCart);
    } 
    else {
      removeCartItem(id);
    }
  };

  const handleSubmit = async () => {
    try {
      const orderDetails = cart.map(cartItem => ({
        item: cartItem.item.name,
        quantity: cartItem.quantity,
        price: cartItem.item.price,
      }));

      await axios.post(
        `${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/createOrder`,
        {
          email: user.Email,
          cart: orderDetails,
          userId: user.id,
          totalPrice: totalPrice,
        },
        {
          headers: { Authorization: `Bearer ${user.Token}` },
        }
      );

      // reset cart after order submission
      setCart([]);
      setTotalPrice(0);
      localStorage.removeItem('cart');
      localStorage.removeItem('totalPrice');

      toastr.success('Your order has successfully been sent for processing!');
    } catch (error) {
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
        await axios.put(
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
      <div className="row row-cols-1 row-cols-md-2 g-4">
        {items?.map((item: Item) => (
          <div key={item.id} className="col">
            <div className="card bg-dark text-light">
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
            {cart.map(cart => (
              <li key={cart.item.id} className="list-group-item d-flex justify-content-between align-items-center">
                {/* Updated cart item layout */}
                <div className="d-flex justify-content-between align-items-center w-50">
                  <span>{cart.item.name}</span>
                  <span className="badge bg-primary rounded-pill mx-2">{cart.quantity}</span>
                  <div>
                    <button className="btn btn-sm btn-secondary me-2" onClick={() => incrementCartItemQty(cart.item.id)}>
                      +
                    </button>
                    <button className="btn btn-sm btn-secondary" onClick={() => decrementCartItemQty(cart.item.id)}>
                      -
                    </button>
                  </div>
                </div>
                <div className="w-50 text-end">
                  ${Number((cart.item.price * cart.quantity).toFixed(2))}
                </div>
              </li>
            ))}
          </ul>
          <div className="mt-3">
            <h5>Cart total: ${totalPrice.toFixed(2)}</h5>
            <button className="btn btn-success" onClick={handleSubmit}>
              Submit order
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default Shop;