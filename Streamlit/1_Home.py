# -*- coding: utf-8 -*-
"""
Spyder Editor

This is a temporary script file.
"""

import streamlit as st
import json
import ipaddress

if 'bool_device_added' not in st.session_state:
    st.session_state.bool_device_added = False
#bool_device_added = False

def string_to_bool(s):
    return s.lower() == 'true'

def is_valid_ip(ip):
    try:
        ipaddress.ip_address(ip)
        return True
    except ValueError:
        return False

def main():
    final_json= {}
    device_data = {}
    register_data= {}


        
    st.set_page_config(
         page_title = "Config Builder",
         page_icon = "😱",
    )
    
                
    
          
    tab1, tab2 = st.tabs(["DeviceInfo", "RegisterInfo"])

    with tab1:
        
        st.header("Deviceinfo")
        option = st.selectbox(
        'Protocol',     
        ('TCP', 'RTU'))
        
        
        if (option == "TCP"):
          slaveID = st.text_input("Slave ID", value=0)
          base = st.selectbox("Base", options=[16,32,64], index=0)
          endianness = st.selectbox("Endianness", options=["Big Endian(true)", "Little Endian(false)"], index=0)
          timeout = st.selectbox("Timeout", options=[1000, 2000, 3000, 4000, 5000], index=0)
          registerSwap = st.selectbox("Register Swap", options=["True", "False"], index=0)
          byteSwap = st.selectbox("Byte Swap", options=["True", "False"], index=1)
          sendFullAddress = st.selectbox("Send Full Address", options=["True", "False"], index=0)
          deviceType = st.text_input("Device Type", value="Controller")
          deviceBrand = st.text_input("Device Brand", value="Bluelog")
          deviceModel = st.text_input("Device Model", value="MC")
          phaseConfig = st.text_input("Phase Config", value="3P")
          ipAddress = st.text_input("IP Address", value="10.6.70.5")
          if not (is_valid_ip(ipAddress)):
            st.error(f"{ipAddress} is not a valid IP address. Please enter a valid IP.")
          port = st.number_input("Port", min_value=0, value=502)
          b_endianness = True
          if endianness == "Big Endian(true)":
              b_endianness = True
          else: b_endianness= False
          device_data = {
             "slaveID": int(slaveID),
             "base": base,
             "endianness": b_endianness,
             "timeout": timeout,
             "registerSwap": string_to_bool(registerSwap),
             "byteSwap" : string_to_bool(byteSwap),
             "sendFullAddress": string_to_bool(sendFullAddress),
             "deviceType":deviceType,
             "deviceBrand": deviceBrand,
             "deviceModel":deviceModel,
             "phaseConfig":phaseConfig,
             "ipAddress": ipAddress,
             "port": port
             
             
          }
        else:
         baudrate = st.selectbox("Baud Rate", options=[9600, 19200, 38400, 57600, 115200], index=0)
         slaveID = st.text_input("Slave ID", value=0)
         base_number = st.selectbox("Base", options=[16,32,64], index=0)
         endianness = st.selectbox("Endianness", options=["Big Endian(true)", "Little Endian(false)"], index=0)
         registerSwap = st.selectbox("Register Swap", options=["True", "False"], index=0)
         byteSwap = st.selectbox("Byte Swap", options=["True", "False"], index=1)
         timeout = st.selectbox("Timeout", options=[1000, 2000, 3000, 4000, 5000], index=0)
         sendFullAddress = st.selectbox("Send Full Address", options=["True", "False"], index=0)
         deviceType = st.selectbox("Device Type", options=["Hybrid", "Digital", "Analog"], index=0)
         deviceBrand = st.text_input("Device Brand", value="Solis")
         deviceModel = st.text_input("Device Model", value="S6")
         phaseConfig = st.text_input("Phase Config", value="1P")
         writeFunctionCode = st.selectbox("Write Function Code", options=[1, 2, 3,4,16], index=2)  
         b_endianness = True
         if endianness == "Big Endian(true)":
             b_endianness = True
         else: b_endianness= False
         device_data = {
            "baud": baudrate,
            "slaveID": int(slaveID),
            "base": base_number,
            "endianness": b_endianness,
            "registerSwap": string_to_bool(registerSwap),
            "byteSwap" : string_to_bool(byteSwap),
            "timeout": timeout,
            "sendFullAddress": string_to_bool(sendFullAddress),
            "deviceBrand": deviceBrand,
            "deviceType":deviceType,
            "deviceModel":deviceModel,
            "phaseConfig":phaseConfig,
            "writeFunctionCode": writeFunctionCode,
            
         }
         
        if st.button("Add device Info"):
           st.session_state.bool_device_added = True 
                
         

    with tab2:
        register_info = []
        st.write("Enter registerInfo data:")
        add_more = True
        counter = 1
        if st.session_state.bool_device_added:
               while add_more:
                    col1, col2, col3, col4, col5 = st.columns(5)

                    with col1:
                       name = st.text_input("Tag", key=f"name_{counter}", value = "invP")

                    with col2:
                        type_ = st.number_input("FunctionCode", value=3, key=f"type_{counter}")

                    with col3:
                        address = st.number_input("Register address", value=0, key=f"address_{counter}")

                    with col4:
                        length = st.number_input("Length of read", value=2, key=f"length_{counter}")

                    with col5:
                        scale = st.number_input("Offset", value=0, key=f"scale_{counter}")
                    col6, col7, col8, col9 = st.columns(4)

                    with col6:
                        register = st.number_input("Scale", value=0, key=f"register_{counter}")

                    with col7:
                        count = st.number_input("DataType", value=4, key=f"count_{counter}")

                    with col8:
                       unit = st.selectbox("Unit", options=["W", "A", "W", "Hz", "%","kWh","none"], index=0)
                       if unit=="none":
                        unit = ""
                    with col9:
                        access = st.number_input("tableType", value=1, key=f"access_{counter}")

                    value = [name, type_, address, length, scale, register, count, unit, access]
                    register_info.append(value)
                    add_more = st.checkbox(f" Add register {counter+1}?")
                    counter +=1
#        if st.button("Add Register Info"):
           
      
 
    file_path = "data.json"
    final_json= {
       "DeviceInfo": device_data,
       "RegisterInfo": register_info
        }
 # Customizing JSON formatting
 # You can specify the indentation level and separators
 # Here, we set indentation to 4 spaces and use a comma and space after each key-value pair
    custom_formatting ={
         "indent": 2,
         "separators": (", ", ": ")   
         
     }
    with open(file_path, "w") as json_file:
        json.dump(final_json, json_file, **custom_formatting) 
    
    
if __name__ == '__main__':
    main()