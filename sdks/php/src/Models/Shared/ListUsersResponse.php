<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


/**
 * ListUsersResponse - List of users
 * 
 * @package formance\stack\Models\Shared
 * @access public
 */
class ListUsersResponse
{
    /**
     * $data
     * 
     * @var ?array<\formance\stack\Models\Shared\User> $data
     */
	#[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('array<formance\stack\Models\Shared\User>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $data = null;
    
	public function __construct()
	{
		$this->data = null;
	}
}