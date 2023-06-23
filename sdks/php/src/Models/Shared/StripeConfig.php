<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class StripeConfig
{
	#[\JMS\Serializer\Annotation\SerializedName('apiKey')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $apiKey;
    
    /**
     * Number of BalanceTransaction to fetch at each polling interval.
     * 
     * 
     * 
     * @var ?int $pageSize
     */
	#[\JMS\Serializer\Annotation\SerializedName('pageSize')]
    #[\JMS\Serializer\Annotation\Type('int')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?int $pageSize = null;
    
    /**
     * The frequency at which the connector will try to fetch new BalanceTransaction objects from Stripe API.
     * 
     * 
     * 
     * @var ?string $pollingPeriod
     */
	#[\JMS\Serializer\Annotation\SerializedName('pollingPeriod')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $pollingPeriod = null;
    
	public function __construct()
	{
		$this->apiKey = "";
		$this->pageSize = null;
		$this->pollingPeriod = null;
	}
}