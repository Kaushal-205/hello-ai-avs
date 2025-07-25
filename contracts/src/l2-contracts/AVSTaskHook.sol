// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {IBN254CertificateVerifierTypes} from
    "@eigenlayer-contracts/src/contracts/interfaces/IBN254CertificateVerifier.sol";
import {OperatorSet} from "@eigenlayer-contracts/src/contracts/libraries/OperatorSetLib.sol";

import {IAVSTaskHook} from "@hourglass-monorepo/src/interfaces/avs/l2/IAVSTaskHook.sol";

contract AVSTaskHook is IAVSTaskHook {
    function validatePreTaskCreation(
        address, /*caller*/
        OperatorSet memory, /*operatorSet*/
        bytes memory /*payload*/
    ) external view {
        //TODO: Implement
    }

    function validatePostTaskCreation(
        bytes32 /*taskHash*/
    ) external {
        //TODO: Implement
    }

    function validateTaskResultSubmission(
        bytes32, /*taskHash*/
        IBN254CertificateVerifierTypes.BN254Certificate memory /*cert*/
    ) external {
        //TODO: Implement
    }
}
