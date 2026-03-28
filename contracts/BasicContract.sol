// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

contract BasicContract {
    uint256 public value;
    event ValueUpdates(uint256 indexed newValue);
    constructor() {
        value = 123456;
    }
    
    function get() public view returns (uint256) {
        return value;
    }
    
    function set(uint256 _value) public {

        value = _value;
        emit ValueUpdates(_value);
    }
}
