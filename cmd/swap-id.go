package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/furya-official/mage/x/bep3/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/furya-official/mgtool/binance"
)

var (
	mageDeputiesStrings map[string]string = map[string]string{
		"bnb":  "mage1r4v2zdhdalfj2ydazallqvrus9fkphmglhn6u6",
		"btcb": "mage14qsmvzprqvhwmgql9fr0u3zv9n2qla8zhnm5pc",
		"busd": "mage1hh4x3a4suu5zyaeauvmv7ypf7w9llwlfufjmuu",
		"xrpb": "mage1c0ju5vnwgpgxnrktfnkccuth9xqc68dcdpzpas",
	}
	bnbDeputiesStrings map[string]string = map[string]string{
		"bnb":  "bnb1jh7uv2rm6339yue8k4mj9406k3509kr4wt5nxn",
		"btcb": "bnb1xz3xqf4p2ygrw9lhp5g5df4ep4nd20vsywnmpr",
		"busd": "bnb10zq89008gmedc6rrwzdfukjk94swynd7dl97w8",
		"xrpb": "bnb15jzuvvg2kf0fka3fl2c8rx0kc3g6wkmvsqhgnh",
	}
)

// SwapIDCmd returns a command to calculate a bep3 swap ID for binance and mage chains.
func SwapIDCmd(cdc *codec.Codec) *cobra.Command {

	mageDeputies := map[string]sdk.AccAddress{}
	for k, v := range mageDeputiesStrings {
		mageDeputies[k] = mustMageAccAddressFromBech32(v)
	}
	bnbDeputies := map[string]binance.AccAddress{}
	for k, v := range bnbDeputiesStrings {
		bnbDeputies[k] = mustBnbAccAddressFromBech32(v)
	}

	cmd := &cobra.Command{
		Use:   "swap-id random_number_hash original_sender_address deputy_addres_or_denom",
		Short: "Calculate binance and mage swap IDs given swap details.",
		Long: fmt.Sprintf(`A swap's ID is: hash(swap.RandomNumberHash, swap.Sender, swap.SenderOtherChain)
One of the senders is always the deputy's address, the other is the user who initiated the first swap (the original sender).
Corresponding swaps on each chain have the same RandomNumberHash, but switched address order.
		
The deputy can be one of %v to automatically use the mainnet deputy addresses, or an arbitrary address.
The original sender and deputy address cannot be from the same chain.
`, getKeys(mageDeputiesStrings)),
		Example: "swap-id 464105c245199d02a4289475b8b231f3f73918b6f0fdad898825186950d46f36 bnb10rr5f8m73rxgnz9afvnfn7fn9pwhfskem5kn0x busd",
		Args:    cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {

			randomNumberHash, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			// try and decode the bech32 address as either mage or bnb
			addressMage, errMage := sdk.AccAddressFromBech32(args[1])
			addressBnb, errBnb := binance.AccAddressFromBech32(args[1])

			// fail if both decoding failed
			isMageAddress := errMage == nil && errBnb != nil
			isBnbAddress := errMage != nil && errBnb == nil
			if !isMageAddress && !isBnbAddress {
				return fmt.Errorf("can't unmarshal original sender address as either mage or bnb: (%s) (%s)", errMage.Error(), errBnb.Error())
			}

			// calculate swap IDs
			depArg := args[2]
			var swapIDMage, swapIDBnb []byte
			if isMageAddress {
				// check sender isn't a deputy
				for _, dep := range mageDeputies {
					if addressMage.Equals(dep) {
						return fmt.Errorf("original sender address cannot be deputy address: %s", dep)
					}
				}
				// pick deputy address
				var bnbDeputy binance.AccAddress
				bnbDeputy, ok := bnbDeputies[depArg]
				if !ok {
					bnbDeputy, err = binance.AccAddressFromBech32(depArg)
					if err != nil {
						return fmt.Errorf("can't unmarshal deputy address as bnb address (%s)", err)
					}
				}
				// calc ids
				swapIDMage = types.CalculateSwapID(randomNumberHash, addressMage, bnbDeputy.String())
				swapIDBnb = binance.CalculateSwapID(randomNumberHash, bnbDeputy, addressMage.String())
			} else {
				// check sender isn't a deputy
				for _, dep := range bnbDeputies {
					if bytes.Equal(addressBnb, dep) {
						return fmt.Errorf("original sender address cannot be deputy address %s", dep)
					}
				}
				// pick deputy address
				var mageDeputy sdk.AccAddress
				mageDeputy, ok := mageDeputies[depArg]
				if !ok {
					mageDeputy, err = sdk.AccAddressFromBech32(depArg)
					if err != nil {
						return fmt.Errorf("can't unmarshal deputy address as mage address (%s)", err)
					}
				}
				// calc ids
				swapIDBnb = binance.CalculateSwapID(randomNumberHash, addressBnb, mageDeputy.String())
				swapIDMage = types.CalculateSwapID(randomNumberHash, mageDeputy, addressBnb.String())
			}

			outString, err := formatResults(swapIDMage, swapIDBnb)
			if err != nil {
				return err
			}
			fmt.Println(outString)
			return nil
		},
	}

	return cmd
}

func formatResults(swapIDMage, swapIDBnb []byte) (string, error) {
	result := struct {
		MageSwapID string `yaml:"mage_swap_id"`
		BnbSwapID  string `yaml:"bnb_swap_id"`
	}{
		MageSwapID: hex.EncodeToString(swapIDMage),
		BnbSwapID:  hex.EncodeToString(swapIDBnb),
	}
	bz, err := yaml.Marshal(result)
	return string(bz), err
}

func mustMageAccAddressFromBech32(address string) sdk.AccAddress {
	a, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return a
}

func mustBnbAccAddressFromBech32(address string) binance.AccAddress {
	a, err := binance.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return a
}

func getKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
